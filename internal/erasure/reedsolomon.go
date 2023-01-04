package erasure

import (
	"cess-portal/conf"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

const (
	Redundancy = 1 * 10 / 15
)

func reedSolomonRule(fsize int64) (int, int) {
	if fsize <= conf.SIZE_1MB*2560 {
		if fsize <= conf.SIZE_1KB {
			return 1, 0
		}

		if fsize <= conf.SIZE_1MB*8 {
			return 2, 1
		}

		if fsize <= conf.SIZE_1MB*64 {
			return 4, 2
		}

		if fsize <= conf.SIZE_1MB*384 {
			return 6, 3
		}

		if fsize <= conf.SIZE_1MB*1024 {
			return 8, 4
		}

		return 10, 5
	}

	if fsize <= conf.SIZE_1MB*6144 {
		return 12, 6
	}

	if fsize <= conf.SIZE_1MB*7168 {
		return 14, 7
	}

	if fsize <= conf.SIZE_1MB*8192 {
		return 16, 8
	}

	if fsize <= conf.SIZE_1MB*9216 {
		return 18, 9
	}

	return 20, 10
}

func ReedSolomon(fpath string, size int64) ([]string, int, int, error) {
	var shardspath = make([]string, 0)
	datashards, rdunshards := reedSolomonRule(size)
	if rdunshards == 0 {
		shardspath = append(shardspath, fpath)
		return shardspath, datashards, rdunshards, nil
	}

	if datashards+rdunshards <= 6 {
		enc, err := reedsolomon.New(datashards, rdunshards)
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}

		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}

		// Split the file into equally sized shards.
		shards, err := enc.Split(b)
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}
		// Encode parity
		err = enc.Encode(shards)
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}
		// Write out the resulting files.
		for i, shard := range shards {
			var outfn = fmt.Sprintf("%s.00%d", fpath, i)
			err = ioutil.WriteFile(outfn, shard, os.ModePerm)
			if err != nil {
				return shardspath, datashards, rdunshards, err
			}
			shardspath = append(shardspath, outfn)
		}
		return shardspath, datashards, rdunshards, nil
	}

	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(datashards, rdunshards)
	if err != nil {
		return shardspath, datashards, rdunshards, err
	}

	f, err := os.Open(fpath)
	if err != nil {
		return shardspath, datashards, rdunshards, err
	}

	instat, err := f.Stat()
	if err != nil {
		return shardspath, datashards, rdunshards, err
	}

	shards := datashards + rdunshards
	out := make([]*os.File, shards)

	// Create the resulting files.
	dir, file := filepath.Split(fpath)

	for i := range out {
		var outfn string
		if i < 10 {
			outfn = fmt.Sprintf("%s.00%d", file, i)
		} else {
			outfn = fmt.Sprintf("%s.0%d", file, i)
		}
		out[i], err = os.Create(filepath.Join(dir, outfn))
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}
		out[i].Close()
		shardspath = append(shardspath, filepath.Join(dir, outfn))
	}

	// Split into files.
	data := make([]io.Writer, datashards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, instat.Size())
	if err != nil {
		return shardspath, datashards, rdunshards, err
	}

	// Close and re-open the files.
	input := make([]io.Reader, datashards)

	for i := range data {
		f, err := os.Open(out[i].Name())
		if err != nil {
			return shardspath, datashards, rdunshards, err
		}
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, rdunshards)
	for i := range parity {
		parity[i] = out[datashards+i]
		defer out[datashards+i].Close()
	}

	// Encode parity
	err = enc.Encode(input, parity)
	if err != nil {
		return shardspath, datashards, rdunshards, err
	}

	return shardspath, datashards, rdunshards, nil
}

func ReedSolomon_Restore(dir, fid string, datashards, rdushards int, fsize uint64) error {
	outfn := filepath.Join(dir, fid)
	_, err := os.Stat(outfn)
	if err == nil {
		return nil
	}
	if datashards+rdushards <= 6 {
		enc, err := reedsolomon.New(datashards, rdushards)
		if err != nil {
			return err
		}
		shards := make([][]byte, datashards+rdushards)
		for i := range shards {
			infn := fmt.Sprintf("%s.00%d", outfn, i)
			shards[i], err = ioutil.ReadFile(infn)
			if err != nil {
				shards[i] = nil
			}
		}

		// Verify the shards
		ok, _ := enc.Verify(shards)
		if !ok {
			err = enc.Reconstruct(shards)
			if err != nil {
				return err
			}
			ok, err = enc.Verify(shards)
			if !ok {
				return err
			}
		}
		f, err := os.Create(outfn)
		if err != nil {
			return err
		}
		defer f.Close()
		err = enc.Join(f, shards, len(shards[0])*datashards)
		return err
	}

	enc, err := reedsolomon.NewStream(datashards, rdushards)
	if err != nil {
		return err
	}

	// Open the inputs
	shards, size, err := openInput(datashards, rdushards, outfn)
	if err != nil {
		return err
	}

	// Verify the shards
	ok, err := enc.Verify(shards)
	if !ok {
		shards, size, err = openInput(datashards, rdushards, outfn)
		if err != nil {
			return err
		}

		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				var outfn string
				if i < 10 {
					outfn = fmt.Sprintf("%s.00%d", outfn, i)
				} else {
					outfn = fmt.Sprintf("%s.0%d", outfn, i)
				}
				out[i], err = os.Create(outfn)
				if err != nil {
					return err
				}
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			return err
		}

		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				if err != nil {
					return err
				}
			}
		}
		shards, size, err = openInput(datashards, rdushards, outfn)
		ok, err = enc.Verify(shards)
		if !ok {
			return err
		}
		if err != nil {
			return err
		}
	}

	f, err := os.Create(outfn)
	log.Println("ReedSolomon_Restore->outfn", outfn)
	if err != nil {
		return err
	}
	defer f.Close()
	shards, size, err = openInput(datashards, rdushards, outfn)
	if err != nil {
		return err
	}

	err = enc.Join(f, shards, int64(datashards)*size)
	return err
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		var infn string
		if i < 10 {
			infn = fmt.Sprintf("%s.00%d", fname, i)
		} else {
			infn = fmt.Sprintf("%s.0%d", fname, i)
		}
		f, err := os.Open(infn)
		if err != nil {
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, 0, err
		}
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}

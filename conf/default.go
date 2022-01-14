package conf

//default set up about cess client
var (
	Conf_File_Path_D = "/etc/cess.d/cess_client.yaml"
	Board_Path_D     = "/etc/cess.d/"
)

/*
	system set up
*/
const (
	Exit_Normal         = 0
	Exit_CmdLineParaErr = -1
	Exit_ConfErr        = -2
)

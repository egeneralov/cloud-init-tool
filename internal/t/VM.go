package t

type VMOptions struct {
	Name    string
	CPUs    int
	Memory  int
	BootISO string
	InitISO string
	Network string
	Args    string
	Disks   []Disk
}

type Disk struct {
	//   Name string
	Size          int
	CreateOptions string
}

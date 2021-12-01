package shortener

//go:generate mockery --name=Backup
type Backup interface {
	Append(line string) error
	ReadAll() map[string]string
}

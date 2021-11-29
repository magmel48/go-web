package shortener

type Backup interface {
	Append(line string)
	ReadAll() map[string]string
}


Using `configer.StringSlice` applications can load multiple files through the command line using the `flag` package.

```bash
$ go run main.go \
  --config config.env \
  --config config.json \
  --config config.yaml
```
```
{Name:loaded from YAML Log:{Level:DEBUG File:/tmp/example.log} Environment:production IDs:[2 4 6 8]}
```

Changing the order will overwrite previous data set in another config file.
```bash
$ go run main.go \
  --config config.yaml \
  --config config.env \
  --config config.json
```
```
{Name:loaded from JSON Log:{Level:DEBUG File:/tmp/example.log} Environment:production IDs:[1 3 5 7]}
```
Notice the difference in `Name` value.

Setting variables in the OS's environment will override any values provided in config files.
```bash
$ NAME="val from env" go run main.go --config config.yaml
```
```
{Name:val from env Log:{Level:DEBUG File:} Environment: IDs:[2 4 6 8]}
```

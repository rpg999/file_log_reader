# Test project

This test code is designed to create a simple log's files tracker.

To get the latest stable version:
```
git clone https://github.com/rpg999/file_log_reader.git
cd file_log_reader
make init
```

### Build
```
make build
```

After that, the new binary will be created in the folder `bin`

### Run
CAUTION: The daemon will try to connect to `mongodb://localhost:27017` and stores the log entries in the` bitsane_test` database, `log` collection

In order to start a test run, you need to create the `/tmp/test_files` folder and copy 2 files
`first_format.log` and` second_format.log`, which are located in the `src/github.com/rpg999/file_log_reader/test_files` to this folder, then call:

```
make test-run
```

Binary file run parameters :
1. `-file_list=/path/to/file` - required
- An example file list can be found in `src/github.com/rpg999/file_log_reader/test_files/test_file_list.json`
- An example of a command can be found in the file `src/github.com/Makefile`, command `make test-run`




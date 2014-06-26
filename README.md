# Shuttle

Shuttle is a minimalisting application deployment tool.

## Compiling

Shuttle is written in Go. You need to Go >= 1.2 to compile Shuttle.

To compile run the following command:

```
make deps
make
```

## Command Reference

```
shuttle setup - Create application structure
shuttle deploy - Run deployment from git repository
shuttle lock - Lock deployment
shuttle unlock - Unlock deployment
shuttle connect - Connect via SSH (not implemented)
shuttle rollback - Rollback to previous release (not implemented)
```

## License

The MIT License (MIT)

Copyright (c) 2014 Dan Sosedoff, <dan.sosedoff@gmail.com>
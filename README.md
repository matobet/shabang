# Shabang

Simple finder of SHA-256 hash collisions written in Golang.

## Dependencies

[Glide](https://github.com/Masterminds/glide) is used for dependency management so please
follow its install instructions and then in the root of this project run `glide install`.

## Run

You can run the project either with

```
$ go run main.go
```

or 

```
$ go build
$ ./shabang --bitlen 24 --seed 'My cool hash seed'
```

# Rabbit Hat
RabbitMQ virtual host cloning tool.

## Usage 

### Rabbit Hat arguments

```
Usage of rabbithat:
  -source-file string
    	File to read source rabbit data (json format)
  -source-password string
    	Source rabbit password
  -source-rabbit string
    	Source rabbit address:port
  -source-user string
    	Source rabbit username
  -source-vhost string
    	Source rabbit virtual host
  -target-file string
    	File to dump source rabbit data (json format)
  -target-password string
    	Target rabbit password
  -target-rabbit string
    	Target rabbit address:port
  -target-user string
    	Target rabbit username
  -target-vhost string
    	Target rabbit virtual host
  -version
    	Rabbit Hat version
```

## Building 

After setting Go environment values 
([goenv.sh](https://github.com/diegohce/nexer/blob/master/goenv.sh) might help), 
go to ```src``` directory and run from the command line:

```go build rabbithat.go```






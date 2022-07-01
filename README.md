# api-playground
[![default](https://github.com/pampatzoglou/api/actions/workflows/default.yaml/badge.svg)](https://github.com/pampatzoglou/api/actions/workflows/default.yaml)

`docker-compose up --detach --build`

http://localhost:8000/metrics
http://localhost:8000/playground

```
{
	__schema {
		queryType {
			fields {
				name
		}
	  }
	}
}
```

http://localhost:8000/query

`docker-compose down`


read: https://gqlgen.com/getting-started/

To get profiling data please visit: `http://localhost:8000/debug/pprof/`

Example for getting heap profile: `curl http://localhost:8000/debug/pprof/heap`
If the profile is needed for last ’n’ seconds you can profile it by setting ‘seconds’ parameter in the query.
`http://localhost:8000/debug/pprof/heap?seconds=n`

You can easily save profile by `curl http://localhost:8000/debug/pprof heap --output heap.tar.gz`

To analyse the profile we can use `go tool pprof <file_created_from_previous_step>`
ie: `go tool pprof heap.tar.gz`

To simplify the Step 2 and Step 3, we can get the profile and display the profile in neat diagram using

`go tool pprof -web http://localhost:8000/debug/pprof/heap`

TODO: https://github.com/ddosify/ddosify

## health
`http://localhost:9000/live`
`http://localhost:9000/ready`

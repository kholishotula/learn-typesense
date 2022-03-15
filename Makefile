typesense:
	docker run -p 8108:8108 -v/tmp/data:/data typesense/typesense:0.22.2 --data-dir /data --api-key=Hu52dwsas2AdxdE

run:
	go run .\cmd\server\main.go
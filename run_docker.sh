docker stop oyvindsk.com
docker rm oyvindsk.com
docker run --name oyvindsk.com -v /root/docker/volume-oyvindsk.com:/app -d -p 82:3001 ubuntu bash -c 'cd /app && /app/goblog 2>&1'

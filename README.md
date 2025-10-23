# vm-backend
Backend for vendor merchant apps




docker run --name merchant-mysql \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=merchantdb \
  -p 3306:3306 \
  -v ~/mysql_data:/var/lib/mysql \
  -d mysql:8

# sooa_user_token_ms
Registro y autenticación basados en tokens para usuarios de SOOA

docker build --tag sooa-user-token .
docker tag sooa-user-token degarzonm/sooa-user-token-ms:v1.1
docker push degarzonm/sooa-user-token-ms:v1.1
docker run -p 6666:6666 sooa-user-token

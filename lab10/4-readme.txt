# Создаём приватный ключ CA                                                      
openssl genrsa -out ca.key 4096

# Создаём самоподписанный сертификат CA
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt \
-subj "/C=RU/ST=SomeState/L=SomeCity/O=MyOrganization/OU=MyUnit/CN=MyRootCA"


# Создаём приватный ключ сервера        
openssl genrsa -out server.key 4096

# Создаём CSR для сервера с использованием конфигурационного файла
openssl req -new -key server.key -out server.csr -config server.cnf

# Подписываем сертификат сервера с помощью CA и включаем расширения
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
-out server.crt -days 3650 -sha256 -extfile server.cnf -extensions req_ext

Certificate request self-signature ok
subject=C=RU, ST=SomeState, L=SomeCity, O=MyOrganization, OU=MyServer, CN=localhost


# Получаем модуль из приватного ключа                                                                
openssl rsa -noout -modulus -in server.key | openssl md5

# Получаем модуль из сертификата
openssl x509 -noout -modulus -in server.crt | openssl md5

MD5(stdin)= 5f51bfcb73e29852c96f42f316950c38
MD5(stdin)= 5f51bfcb73e29852c96f42f316950c38


# Создаём приватный ключ клиента                                                        
openssl genrsa -out client.key 4096

# Создаём CSR для клиента с использованием конфигурационного файла
openssl req -new -key client.key -out client.csr -config client.cnf

# Подписываем сертификат клиента с помощью CA и включаем расширения
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
-out client.crt -days 3650 -sha256 -extfile client.cnf -extensions req_ext

Certificate request self-signature ok
subject=C=RU, ST=SomeState, L=SomeCity, O=MyOrganization, OU=MyClient, CN=client


# Получаем модуль из приватного ключа клиента                            
openssl rsa -noout -modulus -in client.key | openssl md5

# Получаем модуль из сертификата клиента
openssl x509 -noout -modulus -in client.crt | openssl md5

MD5(stdin)= 43970790f317eb0e6a75ef7355424ed7
MD5(stdin)= 43970790f317eb0e6a75ef7355424ed7
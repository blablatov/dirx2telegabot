## Generate private RSA key

To generate RSA key using OpenSSL tool, we need to use `genrsa` command like below.
Генерируем закрытый ключ сертификата:

```shell script
$ openssl genrsa<1> -out server.key<2> 2048<3>
```

1. Specifies which algorithm to use to create the key. OpenSSL supports creating keys with a different algorithm like
 RSA, DSA, and ECDSA. All types are practical for use in all scenarios. For example, for web server keys commonly uses RSA. In our case, we need to generate RSA type key.
2. Specifies the name of the generated key. Can have any name with `.key` as extension.
3. Specifies the size of the key. The default size for RSA keys is only 512 bits, which is not secure because an
 intruder can use brute force to recover your private key. So we use a 2048-bit RSA key which is considered to be secure.

Here can also add a passphrase to the key.  


## Generate CA and self-signed certificates
Let’s create a Certificate Authority and self-signed certificate for our example.   
To generate RSA key using OpenSSL tool, execute the following command.  
Создание корневого приватного ключа:  

```shell script
$ openssl genrsa -aes256 -out ca.key 4096
```

Here we create a new private key with a password for the CA. Now we can create the root CA certificate with a validity of two years using the SHA256 hash algorithm.
Создание самоподписанного корневого сертификата:

```shell script
$ openssl req -new -x509 -sha256 -days 730 -key ca.key -out ca.crt
```

As in the server certificate creation earlier, Certificate generation is an interactive process.   
You can make them blank by entering `.`. But you need to give a name for Common Name as mentioned before.  

So now we created both the private key and self-signed certificate of our Certificate Authority. We can verify the root certificate using below command,
можем проверить корневой сертификат, используя следующую команду:

```shell script
$ openssl x509 -noout -text -in ca.crt
```

We can check the validity period for 02 years and the issuer and subject should both be set to the value passed for “Common Name” because this is a root certificate and it is self-signed.

The next step is to create a server private key and certificate. Unlike the previous section, we need to get the certificate signed by our new Certificate Authority(CA). 


## Generate server certificate
Once we have the server private key, we can proceed to create a Certificate Signing Request (CSR).
Execute the following command to create a certificate signing request.  
Генерируем запрос на сертификат CSR. В этом запросе нужно передать конкретные значения в параметре "subj".  
   
Это должны быть один или несколько параметров CN (Common Name), которые будут прописаны в атрибутах сертификата "Subject" и "Subject Alternative Name".      
Фактически CN - это IP адреса и имена хостов, по которым будет осуществляться доступ к серверу.     
Имена хостов могут быть как в формате доменного имени, так и FQDN:        
   
```shell script
openssl req -new -key server.key -subj "/CN=127.0.0.1/CN=localhost/CN=net-mtls-service" -out server.csr
```

Посмотреть получившийся запрос на сертификат и убедиться в его корректности:  

```shell script
openssl req -text -in server.csr
```

After a CSR is generated, we can sign the request and generate the certificate using our own CA certificate.   
Теперь генерируем сертификат X.509 (server.crt) на 365 дней для сервера и подписываем его приватным ключом CA (ca.key): 

```shell script
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extensions SAN -extfile openssl.cnf
```

Файл openssl.cnf должен содержать список SAN (Subject Alternative Name) идентичный тому, что был в CSR запросе на сертификат:  

```shell script
[SAN]
subjectAltName = @alt_names
[alt_names]
IP.1 = 127.0.0.1
DNS.1 = localhost
DNS.2 = net-mtls-service
```

Посмотреть на получившийся сертификат и убедиться в его корректности:  

```shell script
openssl x509 -text -in ca.crt
```

Now we have created server key(server.key) and server certificate(server.crt).   
We can use them to enable mutual TLS in server side.  


## Generate client key and certificate
Generating the client certificate is very similar to creating the server certificate. We need to execute the following commands to create a private key, create a certificate signing request and create a certificate for client application.
Создать закрытый ключ, создать запрос на подпись сертификата и создать сертификат для клиентского приложения:  
  
```shell script
$ openssl genrsa -out client.key 2048
$ openssl req -new -key client.key -out client.csr
$ openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 2 -out client.crt
```

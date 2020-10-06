# Aplicação com Golang + RabbitMQ

Desenvolvimento de aplicação utilizando Golang + RabbitMQ & Docker.<br>
Exemplo de publish e consumer.

## Tecnologias Utilizadas 🚀
* **[Visual Studio Code](https://code.visualstudio.com/)**
* **[Golang - Visual Studio Code (Extensão)](https://code.visualstudio.com/docs/languages/go)**
* **[Golang](https://golang.org/)**
* **[Docker](https://www.docker.com/)**
* **[RabbitMQ](https://www.rabbitmq.com/)**

## Docker 🐋
1. Para utilizar o RabbitMQ através do Docker devemos criar o container com o seguinte comando: 
```
> docker run --hostname rabbitmq-dev --name rabbitmq-dev -p 5672:5672 -p 15672:15672 -d rabbitmq:3-management
```
1. Após executar o comando acessar a URL `localhost:15672/` e com as credencias, **username:** `guest` e **password:** `guest`
   
## Como Executar? 🔥
1. Para executar o `consumer.go`, entrar na pasta: `/consumers` e executar: 
```
> go run consumer.go
```
2. Para executar o `publish.go`, entrar na pasta: `/publishs` e executar:
```
> go run publish.go
```

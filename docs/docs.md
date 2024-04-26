# Documentação do Rate Limiter

O Rate Limiter implementado em Go é uma ferramenta eficaz para controlar o tráfego de requisições a um serviço web, limitando o número de requisições permitidas por segundo com base no endereço IP do solicitante ou em um token de acesso específico. Este documento fornece uma visão geral de como o rate limiter funciona e como ele pode ser configurado.

## Funcionamento

O rate limiter utiliza um armazenamento externo (por exemplo, Redis) para rastrear o número de requisições feitas por um identificador único (endereço IP ou token de acesso) dentro de uma janela de tempo definida. Quando uma requisição é recebida, o rate limiter verifica se o limite de requisições para aquele identificador já foi alcançado. Se o limite não foi alcançado, a requisição é permitida; caso contrário, é negada.

### Componentes Principais

- **StorageLimiter**: É a estrutura central que implementa a lógica do rate limiter. Ela mantém a configuração do limite de requisições, a janela de tempo e a interface para o armazenamento de dados.
- **MiddlewareRate**: Um middleware HTTP que utiliza o `StorageLimiter` para verificar e limitar as requisições recebidas.

### Chave de Identificação

- Para cada requisição, o rate limiter determina uma chave de identificação baseada no endereço IP do solicitante ou no token de acesso, se fornecido.
- Se um token de acesso é fornecido no cabeçalho `API_KEY`, o rate limiter usa esse token como chave; caso contrário, usa o endereço IP.

### Limitação e Janela de Tempo

- O rate limiter permite configurar um limite de requisições (`limit`) e uma janela de tempo (`window`) durante a qual esse limite se aplica.
- Quando o limite de requisições para uma chave específica é alcançado dentro da janela de tempo, novas requisições com essa chave são negadas até que a janela de tempo seja reiniciada.

## Configuração

Para configurar o rate limiter, você precisa definir o limite de requisições, a janela de tempo e configurar o armazenamento de dados.

### Criando um StorageLimiter

```go
storage := // inicialize seu armazenamento aqui, por exemplo, uma instância Redis
limit := 100 // limite de 100 requisições
window := time.Minute // janela de tempo de 1 minuto

limiter := ratelimiter.NewStorageLimiter(storage, limit, window)
```

### Adicionando o Middleware ao Servidor HTTP

```go
http.Handle("/", MiddlewareRate(limiter)(http.HandlerFunc(handler)))
```

Substitua `handler` pela função que lida com suas requisições HTTP.

### Ajustando Limites por IP e Token

- O limite padrão se aplica a identificações baseadas em endereço IP.
- Para configurar limites específicos para tokens de acesso, você pode ajustar a lógica dentro do `MiddlewareRate` ou gerenciar diferentes instâncias do `StorageLimiter` com configurações distintas para diferentes tokens.

## Considerações

- A eficácia do rate limiter em ambientes de alta carga depende da configuração e do desempenho do armazenamento de dados.
- É importante monitorar e ajustar a configuração do rate limiter com base no padrão de tráfego do seu serviço web para evitar falsos positivos ou negativos na limitação de requisições.

Este documento fornece uma visão geral de como configurar e utilizar o rate limiter em sua aplicação Go. A implementação pode ser ajustada conforme necessário para atender às necessidades específicas do seu serviço web.
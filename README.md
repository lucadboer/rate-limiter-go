# Rate Limiter em Go

## Objetivo

Desenvolver um rate limiter em Go configurável para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Descrição

O rate limiter é uma ferramenta essencial para controlar o tráfego de requisições a um serviço web, protegendo-o contra sobrecarga ou ataques de DDoS. Este projeto implementa um rate limiter em Go que limita as requisições com base em dois critérios: endereço IP e token de acesso.

### Funcionalidades

- **Limitação por Endereço IP**: Restringe o número de requisições de um único endereço IP dentro de um intervalo de tempo definido.
- **Limitação por Token de Acesso**: Permite diferentes limites de requisições baseados em tokens de acesso únicos, com a possibilidade de configurar diferentes tempos de expiração para cada token.

### Configuração

- O rate limiter funciona como um middleware injetado ao servidor web.
- Permite a configuração do número máximo de requisições por segundo.
- Oferece a opção de definir o tempo de bloqueio para um IP ou Token após exceder o limite de requisições.
- As configurações podem ser realizadas via variáveis de ambiente ou em um arquivo `.env`.

### Resposta ao Exceder o Limite

Quando o limite é excedido, o sistema responde com:

- Código HTTP: `429`
- Mensagem: `"You have reached the maximum number of requests or actions allowed within a certain time frame."`

### Persistência

- Utiliza o Redis para armazenar e consultar as informações do limiter.
- Implementa uma "strategy" para permitir a fácil troca do Redis por outro mecanismo de persistência.

## Implementação

### Middleware

O middleware intercepta cada requisição, verifica se o IP ou o token de acesso excedeu o limite de requisições permitidas e, em caso afirmativo, bloqueia a requisição retornando o código HTTP 429.

### Estratégia de Persistência

A lógica de limitação é implementada de forma a ser independente do mecanismo de persistência. Isso é alcançado através de uma interface que pode ser implementada por diferentes backends de armazenamento, como Redis, memória, entre outros.

### Exemplos de Uso

- **Limitação por IP**: Se configurado para permitir 5 requisições por segundo por IP e um IP envia 6 requisições em um segundo, a sexta requisição é bloqueada.
- **Limitação por Token**: Se um token tem um limite configurado de 10 requisições por segundo e envia 11 requisições, a décima primeira é bloqueada.

## Configuração e Testes

- As configurações de limitação são definidas via variáveis de ambiente ou arquivo `.env`.
- Inclui testes automatizados para validar a eficácia e robustez do rate limiter.
- Utiliza Docker/Docker Compose para facilitar os testes da aplicação.

## Como Executar

1. Configure as variáveis de ambiente ou o arquivo `.env` com os limites desejados.
2. Utilize o `docker-compose` para subir o serviço Redis necessário para o rate limiter.
3. Inicie o servidor web que responderá na porta 8080.

## Conclusão

Este rate limiter é uma solução eficaz para controlar o tráfego de requisições, protegendo serviços web contra sobrecargas e garantindo uma distribuição justa dos recursos entre os usuários.

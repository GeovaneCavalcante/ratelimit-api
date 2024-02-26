
# Ratelimit

## Descrição
Ratelimit é uma aplicação em Go, projetada para limitar o número de solicitações a um endpoint com base em tokens ou endereços IP. Utiliza Redis para armazenamento de dados e Locust para testes de carga. A aplicação é ideal para gerenciar o tráfego de entrada, garantindo a disponibilidade e a eficiência do serviço.

## Endpoint
A aplicação possui um único endpoint:

```bash
curl --location 'http://localhost:8080/health' \
--header 'API_KEY: 5095bc00-2f9e-4e6f-b355-11688d20530d' \
--header 'X-Forwarded-For: 192.168.0.1'
```
O header `X-Forwarded-For` é opcional. Caso não fornecido, o endereço IP real da máquina será utilizado.

## Configuração
No arquivo `.env`, as seguintes configurações devem ser definidas:

- `TOKENS_CONFIG_LIMIT`: Define limites de requisições e bloqueio por token.
- `IP_CONFIG_LIMIT`: Define limites de requisições e bloqueio por IP.

Exemplo:

```env
TOKENS_CONFIG_LIMIT=[{"token": "5095bc00-2f9e-4e6f-b355-11688d20530d", "max_requests": 30, "block_time_seconds": 60}, {"token": "eeec68b2-f1b9-4adc-813a-4cbade5d5387", "max_requests": 5, "block_time_seconds": 1800}]
IP_CONFIG_LIMIT={"max_requests": 30, "block_time_seconds": 60}
```
### Alterar persistência 
O rate limiter utiliza redis como storage e que permite viabilizar uma `stragegy` que empilha eventos e com base nos mesmo é implementado a regra de negócio com base nas políticas de acesso. Caso queira trocar a persistência e utilizar outra ferramenta é necessário fazer a implementação da interface `EventStorageInterface` que está contida no diretório `pkg/ratelimit/event.go`. 

O pacote `ratelimit` depende de uma implementação dessa interface com regras de persistência para empilhar eventos. A implementação deve ser passada na inicialização no arquivo `cmd/server.go`.



## Construção e Execução
Para construir a imagem da aplicação, utilize:
```bash
make docker-build-image
```
Para subir as dependências e executar a aplicação:
```bash
make docker-up
```
Isso irá iniciar o Redis e o Locust para testes de carga.
## Testes Unitários e Cobertura de Testes

Para executar os testes unitários da aplicação, utilize o comando:

```bash
make test
```

Este comando executa todos os testes unitários disponíveis no projeto, mostrando a saída de forma detalhada.

Para avaliar a cobertura dos testes, execute:

```bash
make test-coverage
```

Este comando gera um relatório de cobertura dos testes (`coverage.out`) e uma representação visual em HTML (`coverage.html`). É útil para identificar partes do código não cobertas pelos testes.

## Testes de Carga com Locust
Endereço do Locust: `http://localhost:8089`

### Resultados de Testes Anteriores
Localizados em: `stress_test/results`

Configurações do teste:

Caso 1

- 10 usuários
- Spawn rate: 1 usuário/segundo
- Duração do teste: 3 minutos
- Token definido e usado nas requisições: `5095bc00-2f9e-4e6f-b355-11688d20530d`
- Configurado a permissão 30 requisições em um intervalo de 1 minuto para token.


Resultados caso 1:
- Número de requests: 117014
- Requests com falha: 116924
- Requests com sucesso: 90


Caso 2
- 10 usuários
- Spawn rate: 1 usuário/segundo
- Duração do teste: 3 minutos
- IP definido e usado nas requisições: `192.168.9.9`
- Configurado a permissão 20 requisições em um intervalo de 1 minuto para token.

Resultados caso 2:
- Número de requests: 117439
- Requests com falha: 117379
- Requests com sucesso: 60


### Executando Novos Testes
Para realizar novos testes de carga, siga estas etapas:
1. Acesse o Locust através do navegador na porta `8089`.
2. Configure o número de usuários e o spawn rate.
3. Inicie o teste e monitore os resultados em tempo real.

## Conclusão
A aplicação Ratelimit demonstrou eficiência ao lidar com altas cargas, mantendo a estabilidade e a funcionalidade do serviço mesmo sob demanda intensa, conforme demonstrado pelos testes anteriores.


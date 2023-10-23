# Chat service - EP02


## Geral
Esse trabalho funciona com 3 arquivos pilares:
* server - arquivo de host do servidor: permite comunicação entre clientes e estabelece alguns comandos.
* client - arquivo de cliente do usuario: sincroniza o input e output do terminal com um canal no servidor.
* bot - modelo de arquivo do bot: utiliza a mesma lógica de conexão do client, mas não envia e nem recebe mensagens de escopo global.


## Comandos

Comando | Descrição
:--: | :--:
/changenickname [nome] | Muda o nickname associado ao ip do client
/checkip (ou /ip) | Recebe uma mensagem privada do servidor informando o ip do client
/msg [destinatário] [mensagem] | Envia uma mensagem privada para o destinatário
/quit (ou /q) | Encerra comunicação do client com servidor

## Bot

O bot utilizado como exemplo é um bot que simplesmente inverte a mensagem privada e manda de volta ao remetente. Ao conectar no servidor ele executa o comando '/changenickname inversor' para mudar seu nickname e facilitar a comunicação com outros usuários.

Exemplo:

Mensagem enviada ao bot -> mensagem retornada pelo bot
``` /msg inversor oioioi``` -> ``` (private) inversor: ioioio```
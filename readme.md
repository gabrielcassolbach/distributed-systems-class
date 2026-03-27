

desenvolver um sistema distribuído baseado em microsserviços para gerenciamento e divulgação de promoções de produtos. O sistema deve seguir uma arquitetura orientada a eventos (Event-drive architecture), na qual os microsserviços se comunicam exclusivamente através de eventos publicados e consumidos em/de um broker RabbitMQ. Cada microsserviço deverá atuar de forma independente e desacoplada. 


Usuários podem cadastrar promoções, votar em promoções cadastradas e receber notificações sobre promoções de seus interesse. 


Promoções com grande quantidade de votos positivos devem ser destacadas como promoções em destaque. 



o sistema deve ser implementado utilizando 4 microsserviços independentes, cada
um responsável por uma parte da lógica da aplicação. 


Microsserviço Gateway: realiza a interação com os usuários (clientes e lojas)
por meio do terminal.


Microsserviço Promocao: responsável pelo gerenciamento das promoções no sistema.
Esse serviço indica gl
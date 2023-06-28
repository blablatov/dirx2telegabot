## Demo packet of chatbot Telegram for Directum RX  

### Описание. Description  
Демонстрационный модуль передачи `rest`-данных из СЭД `Directum RX` в чат-бот `Telegram`.  
При получении определенных данных, можно перейти по гиперссылке, через клик на кнопке клавиатуры чата, 
к нужному документу в Проводнике `Directum RX`.  

В демо пакете токен авторизации для одного пользователя хранится явно в конфиге.    
В полной версии токены пользователей хранятся в БД.  

В `Directum RX` создано демо Решение включающее справочник с фиксированными rest-запррсами,  
при переходе по такой гиперссылке отправляется запрос в модуль обмена с чат-ботом `Telegram`.

### Использование. How use  
По клику на сформированной демо-ссылке в `Directum RX`, отправляется `rest`-сообщение к `http-серверу` чат-бота.  
Пользователь может в любой момент ввести сообщение запроса данных `dirx` в своем чат-боте `Telegram`.  
При вводе ключа `dirx` выводится демо клавиатура чат-бота для взаимодействия с СЭД `Directum RX`. 
Если данных из `Directum RX` нет, видим сообщение `Очередь Directum RX пуста`.  
При вводе произвольных данных, видим сообщение `Введите: dirx"`.  
    
### Блок-схема обмена данными. Block diagram of work.    

			
```mermaid
graph TB

  SubGraph1Flow
  subgraph "ChatBot Telegram"
  SubGraph1Flow(Queue of messages)
  SubGraph1Flow -- update.Message.Chat.ID --> Any_not_system_data --> Notion -->return
  SubGraph1Flow -- update.Message.Chat.ID --> System_message`dirx`--> Numeric_Keyboard --> Node2
  SubGraph1Flow -- update.Message.Chat.ID --> If_empty_queue --> Notion --> return
  end
 
  subgraph "Directum RX"
  SubGraph2Flow(User Folders)
  end

  subgraph "Directum RX"
  Node1[Method with the demo URL] -- Transition to link --> SubGraph1Flow
  Node2[Link to Directum RX User Folders] -- Transition to link --> SubGraph2Flow
end
```   
 

### US  
Before all we should do in `Directum RX` demo URL to `http`-server our the chatbot.  
Demo packet exchange of rest-data from `Directum RX` to chat-bot `Telegram`.   
After got data we can click the desired link keyboard of chat.  
 
  




 

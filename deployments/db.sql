CREATE TABLE "Chats" ( `id` INTEGER, `title` TEXT, `type` TEXT, PRIMARY KEY(`id`) );
CREATE TABLE "Items" ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `category` TEXT, `url` TEXT NOT NULL, `hash` INT UNIQUE, `caption` TEXT );
CREATE TABLE "Views" ( `itemId` INTEGER, `chatId` INTEGER, `requestDate` INTEGER, PRIMARY KEY(`chatId`,`itemId`), FOREIGN KEY(`chatId`) REFERENCES "Chats"(`id`) );

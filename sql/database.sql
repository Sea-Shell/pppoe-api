CREATE TABLE IF NOT EXISTS rating (
    rateId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT DEFAULT 0,
    rateCategoryId INTEGER NOT NULL,
    eventId INTEGER NOT NULL,
    rateName TEXT NOT NULL,
    rateScaleId INTEGER NOT NULL,
    rateScaleValues TEXT NOT NULL,
    FOREIGN KEY (rateCategoryId) REFERENCES rateCategory(categoryId),
    FOREIGN KEY (eventId) REFERENCES event(eventId)
    FOREIGN KEY (rateScaleId) REFERENCES ratingScale(scaleId)
);

CREATE TABLE IF NOT EXISTS rateCategory (
    categoryId INTEGER PRIMARY KEY AUTOINCREMENT DEFAULT 0,
    categoryEventId INTEGER NOT NULL,
    categoryName TEXT NOT NULL,
    FOREIGN KEY (categoryEventId) REFERENCES event(eventId)
);

CREATE TABLE IF NOT EXISTS ratingScale (
    scaleId INTEGER PRIMARY KEY AUTOINCREMENT DEFAULT 0,
    scaleName TEXT NOT NULL,
    scaleDescription TEXT NOT NULL,
    scaleExample TEXT NOT NULL,
    scalePattern TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
    eventId INTEGER PRIMARY KEY AUTOINCREMENT DEFAULT 0,
    eventName TEXT NOT NULL,
    eventCreated DATETIME DEFAULT CURRENT_TIMESTAMP,
    eventStart  DATETIME DEFAULT CURRENT_TIMESTAMP,
    eventEnd    DATETIME DEFAULT CURRENT_TIMESTAMP,
    eventLocation TEXT NOT NULL DEFAULT "None",
    eventDescription TEXT NOT NULL DEFAULT "None",
    eventOrganizer TEXT NOT NULL DEFAULT "None"
);

INSERT INTO event (eventId, eventName) VALUES 
    (1, "PølseBonanza"),
    (2, "HamburgerHumbug"),
    (3, "KebabKaos"),
    (4, "PizzaParty"),
    (5, "TacoTerror"),
    (6, "SushiSlaughter"),
    (7, "PastaPanic"),
    (8, "SaladSlaughter");

INSERT INTO ratingScale (scaleId, scaleName, scaleDescription, scaleExample, scalePattern) VALUES
    (1, "numeric", "A numeric scale", "1-5", "^[1-5]$"),
    (2, "ordinal", "An ordinal scale", "lite, middels, mye, veldig mye, ekstremt mye", "^(lite|middels|mye|veldig mye|ekstremt mye)$"),
    (3, "binary", "A binary scale", "Ja, Nei", "^(Ja|Nei)$");

INSERT INTO rateCategory (categoryId, categoryEventId, categoryName) VALUES
    (1, 1, "Smak"),
    (2, 1, "Tekstur"),
    (3, 1, "Utseende"),
    (4, 1, "Størrelse"),
    (5, 1, "Tilbehør"),
    (6, 1, "Service"),
    (7, 1, "Atmosfære"),
    (8, 1, "Total"),
    (9, 2, "Smak"),
    (10, 2, "Utseende"),
    (11, 2, "Pris"),
    (12, 2, "Størrelse"),
    (13, 2, "Tilbehør"),
    (14, 2, "Service"),
    (15, 2, "Atmosfære"),
    (16, 2, "Total"),
    (17, 3, "Smak"),
    (18, 3, "Utseende"),
    (19, 3, "Pris"),
    (20, 3, "Størrelse"),
    (21, 3, "Tilbehør"),
    (22, 3, "Service"),
    (23, 3, "Atmosfære"),
    (24, 3, "Total"),
    (25, 4, "Smak"),
    (26, 4, "Utseende"),
    (27, 4, "Pris"),
    (28, 4, "Størrelse"),
    (29, 4, "Tilbehør"),
    (30, 4, "Service"),
    (31, 4, "Atmosfære"),
    (32, 4, "Total"),
    (33, 5, "Smak"),
    (34, 5, "Utseende"),
    (35, 5, "Pris"),
    (36, 5, "Størrelse"),
    (37, 5, "Tilbehør"),
    (38, 5, "Service"),
    (39, 5, "Atmosfære"),
    (40, 5, "Total"),
    (41, 6, "Smak"),
    (42, 6, "Utseende"),
    (43, 6, "Pris"),
    (44, 6, "Størrelse"),
    (45, 6, "Tilbehør"),
    (46, 6, "Service"),
    (47, 6, "Atmosfære"),
    (48, 6, "Total"),
    (49, 7, "Smak"),
    (50, 7, "Utseende"),
    (51, 7, "Pris"),
    (52, 7, "Størrelse"),
    (53, 7, "Tilbehør"),
    (54, 7, "Service"),
    (55, 7, "Atmosfære"),
    (56, 7, "Total"),
    (57, 8, "Smak"),
    (58, 8, "Utseende"),
    (59, 8, "Pris"),
    (60, 8, "Størrelse"),
    (61, 8, "Tilbehør"),
    (62, 8, "Service"),
    (63, 8, "Atmosfære"),
    (64, 8, "Total");

INSERT INTO rating (rateId, rateCategoryId, eventId, rateName, rateScaleId, rateScaleValues) VALUES 
    (1, 1, 1, "kjøttsmak", "ordinal", "lite, middels, mye, veldig mye, ekstremt mye"),
    (2, 1, 1, "smaksrikdom", "ordinal", "lav, middels, høy, veldig høy, sjukt nais"),
    (3, 1, 1, "Ettersmak", "ordinal", "lite, middels, mye, veldig mye, ekstremt mye"),
    (4, 2, 1, "Fasthet", "ordinal", "slapp, fast, veldig fast, gummi"),
    (5, 2, 1, "Saftighet", "numeric", "1-4"),
    (6, 2, 1, "Tyggemotstand", "ordinal", "mør, middels, seig, gummi"),
    (7, 3, 1, "Farge", "ordinal", "blass, ok, fin"),
    (8, 3, 1, "Form", "ordinal", "rar, normal, fin"),
    (9, 3, 1, "Størrelse", "ordinal", "liten, stor"),
    (10, 4, 2, "Størrelse", "numeric", "1-3"),
    (11, 4, 2, "Pris", "numeric", "1-3"),
    (12, 4, 2, "Mettende", "numeric", "1-3"),
    (13, 5, 2, "Brød", "numeric", "1-3"),
    (14, 5, 2, "Drikke", "numeric", "1-3"),
    (15, 5, 2, "Tilbehør", "numeric", "1-3"),
    (16, 6, 2, "Ventetid", "numeric", "1-3"),
    (17, 6, 3, "Service", "numeric", "1-3"),
    (18, 6, 3, "Vennlighet", "numeric", "1-3"),
    (19, 7, 3, "Lys", "numeric", "1-3"),
    (20, 7, 3, "Støy", "numeric", "1-3"),
    (21, 7, 3, "Renhold", "numeric", "1-3"),
    (22, 8, 3, "Total", "numeric", "1-5"),
    (23, 9, 4, "Smak", "numeric", "1-5"),
    (24, 9, 4, "Krydder", "numeric", "1-5"),
    (25, 9, 4, "Kjøtt", "numeric", "1-5"),
    (26, 10, 4, "Farge", "numeric", "1-3"),
    (27, 10, 4, "Form", "numeric", "1-3"),
    (28, 10, 4, "Størrelse", "numeric", "1-3"),
    (29, 11, 5, "Pris", "numeric", "1-3"),
    (30, 11, 5, "Mettende", "numeric", "1-3"),
    (31, 11, 5, "Pris", "numeric", "1-3");

/* INSERT INTO users (userUsername, userPassword, userName, userEmail) VALUES 
    ("Bateau", "$2a$10$X2BAOJFWXxAudCm9ShaHvucsdv1.dz3pdbBPf6bJerWs7YJB7KV9", "Mats Bøe Bergmann", "mats.bergmann@gmail.com"); */

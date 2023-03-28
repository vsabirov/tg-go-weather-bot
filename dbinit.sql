CREATE TABLE IF NOT EXISTS UserStats (
  id BIGINT UNIQUE PRIMARY KEY NOT NULL,

  firstRequest TIMESTAMP DEFAULT NOW() NOT NULL,
  amountOfRequests INT DEFAULT 0 NOT NULL,
  lastCity VARCHAR(255) DEFAULT 'Неизвестно' NOT NULL
);

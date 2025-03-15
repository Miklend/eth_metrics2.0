CREATE TABLE IF NOT EXISTS gas_metrics (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Время создания записи
    lastBlock TEXT NOT NULL,  -- Номер последнего блока
    safeGasPrice TEXT NOT NULL,  -- Цена газа для безопасной транзакции
    proposeGasPrice TEXT NOT NULL,  -- Цена газа для предложенной транзакции
    fastGasPrice TEXT NOT NULL,  -- Цена газа для быстрой транзакции
    suggestBaseFee TEXT NOT NULL,  -- Предложенная базовая плата
    PRIMARY KEY (created_at, lastBlock)  -- Можно использовать комбинированный ключ (время и номер блока)
);

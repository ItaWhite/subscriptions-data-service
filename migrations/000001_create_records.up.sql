CREATE TABLE IF NOT EXISTS records (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
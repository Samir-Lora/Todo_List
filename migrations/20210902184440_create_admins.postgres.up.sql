-- Insert rows into table 'TableName'
INSERT INTO users
( -- columns to insert data into
 id, email, password_hash, name, lastname, active, rol, created_at, updated_at
)
VALUES
( -- first row: values for the columns in the list above
 '99d415a2-ea98-44e2-aaf4-6cc8064f5357','samirlora0@gmail.com','$2a$10$eJ65uhrvF7bmoaRTUZmvg.AUcOaehUhHZdHtHVtYHzg1iWmq02llW'  ,'Samir','Lora','active','admin',NOW(), NOW()
);

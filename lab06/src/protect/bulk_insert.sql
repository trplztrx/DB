\copy items(name, description, quantity, price, created_at, updated_at) FROM 'items_data.csv' DELIMITER ',' CSV HEADER;

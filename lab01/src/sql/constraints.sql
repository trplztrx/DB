ALTER TABLE orders
    ADD CONSTRAINT pk_order_id PRIMARY KEY (id),
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN delivery_date SET NOT NULL,
    ALTER COLUMN status SET NOT NULL,
    ALTER COLUMN delivery_cost SET NOT NULL,
    ADD CONSTRAINT chk_order_status CHECK (status IN ('создан', 'отправлен', 'доставлен', 'отменен')),
    ADD CONSTRAINT chk_delivery_cost CHECK (delivery_cost > 0),
    ADD CONSTRAINT fk_sender_user FOREIGN KEY (sender_user_id) REFERENCES users(id),
    ADD CONSTRAINT fk_receiver_user FOREIGN KEY (receiver_user_id) REFERENCES users(id),
    ADD CONSTRAINT fk_pickup_address FOREIGN KEY (pickup_address_id) REFERENCES addresses(id),
    ADD CONSTRAINT fk_delivery_address FOREIGN KEY (delivery_address_id) REFERENCES addresses(id);
    
ALTER TABLE orders
    ALTER COLUMN created_at DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE users
    ADD CONSTRAINT pk_user_id PRIMARY KEY (id),
    ALTER COLUMN user_type SET NOT NULL,
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN phone SET NOT NULL,
    ALTER COLUMN email SET NOT NULL,
    ADD CONSTRAINT chk_user_type CHECK (user_type IN ('физик', 'юрик')),
    ADD CONSTRAINT unique_phone UNIQUE (phone),
    ADD CONSTRAINT unique_email UNIQUE (email);

ALTER TABLE addresses
    ADD CONSTRAINT pk_address_id PRIMARY KEY (id),
    ALTER COLUMN region SET NOT NULL,
    ALTER COLUMN city SET NOT NULL,
    ALTER COLUMN street SET NOT NULL,
    ALTER COLUMN house SET NOT NULL,
    ALTER COLUMN address_type SET NOT NULL,
    ADD CONSTRAINT chk_address_type CHECK (house > 0),
    ADD CONSTRAINT chk_address_type CHECK (address_type IN ('отправитель', 'получатель'));

ALTER TABLE couriers
    ADD CONSTRAINT pk_courier_id PRIMARY KEY (id),
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN surname SET NOT NULL,
    ALTER COLUMN phone SET NOT NULL,
    ALTER COLUMN status SET NOT NULL,
    ADD CONSTRAINT unique_phone UNIQUE (phone),
    ADD CONSTRAINT chk_courier_status CHECK (status IN ('свободен', 'в пути', 'занят'));

ALTER TABLE orderItem
    ADD CONSTRAINT pk_order_item_id PRIMARY KEY (id),
    ALTER COLUMN order_id SET NOT NULL,
    ALTER COLUMN item_name SET NOT NULL,
    ALTER COLUMN quantity SET NOT NULL,
    ADD CONSTRAINT chk_quantity CHECK (quantity > 0),
    ADD CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE orderStatus
    ADD CONSTRAINT pk_order_status_id PRIMARY KEY (id),
    ALTER COLUMN order_id SET NOT NULL,
    ALTER COLUMN status SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD CONSTRAINT chk_order_status CHECK (status IN ('создан', 'отправлен', 'доставлен', 'отменен')),
    ADD CONSTRAINT fk_order_status_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE orderCourier
    ADD CONSTRAINT pk_order_courier PRIMARY KEY (courier_id, order_id),
    ALTER COLUMN courier_id SET NOT NULL,
    ALTER COLUMN order_id SET NOT NULL,
    ALTER COLUMN assigned_at SET NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD CONSTRAINT fk_order_courier_courier FOREIGN KEY (courier_id) REFERENCES couriers(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_order_courier_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;


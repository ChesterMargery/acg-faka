-- SQLite schema for ACG-FAKA
-- Converted from MySQL Install.sql

PRAGMA foreign_keys = ON;

-- Bills table
CREATE TABLE bills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner INTEGER NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    balance DECIMAL(14, 2) NOT NULL,
    type INTEGER NOT NULL, -- 0=expense, 1=income
    currency INTEGER NOT NULL DEFAULT 0, -- 0=balance, 1=coin
    log VARCHAR(64) NOT NULL,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bills_owner ON bills(owner);
CREATE INDEX idx_bills_type ON bills(type);

-- Business table
CREATE TABLE businesses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    shop_name VARCHAR(32),
    title VARCHAR(32),
    notice TEXT,
    service_qq VARCHAR(16),
    service_url VARCHAR(255),
    subdomain VARCHAR(64),
    topdomain VARCHAR(64),
    master_display INTEGER NOT NULL DEFAULT 0, -- 0=no, 1=yes
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_businesses_user_id ON businesses(user_id);
CREATE UNIQUE INDEX idx_businesses_subdomain ON businesses(subdomain);
CREATE UNIQUE INDEX idx_businesses_topdomain ON businesses(topdomain);

-- Business levels table
CREATE TABLE business_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32) NOT NULL,
    icon VARCHAR(255),
    cost DECIMAL(4, 2) NOT NULL DEFAULT 0.00, -- commission percentage
    accrual DECIMAL(4, 2) NOT NULL DEFAULT 0.00, -- profit percentage from main site
    substation INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    top_domain INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- purchase price
    supplier INTEGER NOT NULL DEFAULT 1 -- 0=disabled, 1=enabled
);

-- Insert default business levels
INSERT INTO business_levels (id, name, icon, cost, accrual, substation, top_domain, price, supplier) VALUES
(1, '体验版', '/assets/static/images/business/v1.png', 0.30, 0.10, 1, 0, 188.00, 1),
(2, '普通版', '/assets/static/images/business/v2.png', 0.25, 0.15, 1, 0, 288.00, 1),
(3, '专业版', '/assets/static/images/business/v3.png', 0.20, 0.20, 1, 1, 388.00, 1);

-- Cards table
CREATE TABLE cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner INTEGER NOT NULL DEFAULT 0, -- 0=system, other=user_id
    commodity_id INTEGER NOT NULL,
    draft VARCHAR(255), -- pre-selection info
    secret VARCHAR(760) NOT NULL, -- card/key info
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    purchase_time DATETIME,
    order_id INTEGER,
    status INTEGER NOT NULL DEFAULT 0, -- 0=unsold, 1=sold, 2=locked
    note VARCHAR(64),
    race VARCHAR(32) -- product category
);

CREATE INDEX idx_cards_owner ON cards(owner);
CREATE INDEX idx_cards_commodity_id ON cards(commodity_id);
CREATE INDEX idx_cards_order_id ON cards(order_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_cards_race ON cards(race);

-- Cash withdrawals table
CREATE TABLE cash_withdrawals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    amount DECIMAL(14, 2) NOT NULL,
    type INTEGER NOT NULL DEFAULT 0, -- 0=auto, 1=manual
    card INTEGER NOT NULL, -- 0=alipay, 1=wechat
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    arrive_time DATETIME,
    cost DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- transaction fee
    status INTEGER NOT NULL, -- 0=processing, 1=success, 2=failed, 3=frozen
    message VARCHAR(64)
);

CREATE INDEX idx_cash_withdrawals_user_id ON cash_withdrawals(user_id);
CREATE INDEX idx_cash_withdrawals_type ON cash_withdrawals(type);

-- Categories table
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    sort INTEGER NOT NULL DEFAULT 0,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    owner INTEGER NOT NULL DEFAULT 0, -- 0=system, other=user_id
    icon VARCHAR(255),
    status INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    hide INTEGER NOT NULL DEFAULT 0, -- 0=visible, 1=hidden
    user_level_config TEXT
);

CREATE INDEX idx_categories_owner ON categories(owner);
CREATE INDEX idx_categories_sort ON categories(sort);

-- Insert default category
INSERT INTO categories (id, name, sort, owner, icon, status, hide) VALUES
(1, 'DEMO', 1, 0, '/favicon.ico', 1, 0);

-- Commodities table
CREATE TABLE commodities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cover VARCHAR(255),
    factory_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- guest price
    user_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- member price
    status INTEGER NOT NULL DEFAULT 0, -- 0=offline, 1=online
    owner INTEGER NOT NULL DEFAULT 0, -- 0=system, other=user_id
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    api_status INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    code VARCHAR(64) NOT NULL, -- product code for API
    delivery_way INTEGER NOT NULL DEFAULT 0, -- 0=auto, 1=manual/plugin
    delivery_auto_mode INTEGER NOT NULL DEFAULT 0, -- 0=old_first, 1=random, 2=new_first
    delivery_message VARCHAR(255),
    contact_type INTEGER NOT NULL DEFAULT 0, -- 0=any, 1=phone, 2=email, 3=qq
    password_status INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    sort INTEGER NOT NULL DEFAULT 0,
    coupon INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    shared_id INTEGER,
    shared_code VARCHAR(64),
    shared_premium DECIMAL(10, 2) DEFAULT 0.00,
    shared_premium_type INTEGER DEFAULT 0,
    seckill_status INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    seckill_start_time DATETIME,
    seckill_end_time DATETIME,
    draft_status INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled
    draft_premium DECIMAL(10, 2) DEFAULT 0.00,
    inventory_hidden INTEGER NOT NULL DEFAULT 0, -- 0=visible, 1=hidden
    leave_message VARCHAR(255),
    recommend INTEGER DEFAULT 0, -- 0=no, 1=yes
    send_email INTEGER NOT NULL DEFAULT 0, -- 0=no, 1=yes
    only_user INTEGER NOT NULL DEFAULT 0, -- 0=no, 1=yes (login required)
    purchase_count INTEGER NOT NULL DEFAULT 0, -- 0=unlimited
    widget TEXT,
    level_price TEXT,
    level_disable INTEGER NOT NULL DEFAULT 0,
    minimum INTEGER NOT NULL DEFAULT 0, -- minimum purchase quantity
    maximum INTEGER NOT NULL DEFAULT 0, -- maximum purchase quantity
    shared_sync INTEGER NOT NULL DEFAULT 0,
    config TEXT,
    hide INTEGER NOT NULL DEFAULT 0, -- 0=visible, 1=hidden
    inventory_sync INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_commodities_code ON commodities(code);
CREATE INDEX idx_commodities_owner ON commodities(owner);
CREATE INDEX idx_commodities_status ON commodities(status);
CREATE INDEX idx_commodities_sort ON commodities(sort);
CREATE INDEX idx_commodities_category_id ON commodities(category_id);
CREATE INDEX idx_commodities_api_status ON commodities(api_status);
CREATE INDEX idx_commodities_recommend ON commodities(recommend);

-- Insert demo commodity
INSERT INTO commodities (id, category_id, name, description, cover, factory_price, price, user_price, status, owner, code, delivery_way, sort, coupon, api_status)
VALUES (1, 1, 'DEMO', '<p>该商品是演示商品</p>', '/favicon.ico', 0.00, 1.00, 0.90, 1, 0, '8AE80574F3CA98BE', 1, 1, 1, 1);

-- Configs table
CREATE TABLE configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(128) NOT NULL UNIQUE,
    value TEXT NOT NULL
);

-- Insert default configs
INSERT INTO configs (key, value) VALUES
('shop_name', '异次元店铺'),
('title', '异次元店铺 - 最适合你的个人店铺系统！'),
('description', ''),
('keywords', ''),
('user_theme', 'Cartoon'),
('registered_state', '1'),
('registered_type', '0'),
('registered_verification', '1'),
('registered_phone_verification', '0'),
('registered_email_verification', '0'),
('sms_config', '{"accessKeyId":"","accessKeySecret":"","signName":"","templateCode":""}'),
('email_config', '{"smtp":"","port":"","username":"","password":""}'),
('login_verification', '1'),
('forget_type', '0'),
('notice', '<p><b><font color="#f9963b">本程序为开源程序，使用者造成的一切法律后果与作者无关。</font></b></p>'),
('trade_verification', '1'),
('recharge_welfare', '0'),
('recharge_welfare_config', ''),
('promote_rebate_v1', '0.1'),
('promote_rebate_v2', '0.2'),
('promote_rebate_v3', '0.3'),
('substation_display', '1'),
('domain', ''),
('service_qq', ''),
('service_url', ''),
('cash_type_alipay', '1'),
('cash_type_wechat', '1'),
('cash_cost', '5'),
('cash_min', '100'),
('cname', ''),
('background_url', '/assets/admin/images/login/bg.jpg'),
('default_category', '0'),
('substation_display_list', '[]'),
('closed', '0'),
('closed_message', '我们正在升级，请耐心等待完成。'),
('recharge_min', '10'),
('recharge_max', '1000'),
('user_mobile_theme', '0'),
('commodity_recommend', '0'),
('commodity_name', '推荐'),
('background_mobile_url', ''),
('username_len', '6'),
('cash_type_balance', '0'),
('callback_domain', '');

-- Coupons table
CREATE TABLE coupons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR(32) NOT NULL UNIQUE,
    commodity_id INTEGER NOT NULL,
    owner INTEGER NOT NULL DEFAULT 0,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_time DATETIME,
    service_time DATETIME,
    money DECIMAL(10, 2) NOT NULL,
    status INTEGER NOT NULL DEFAULT 0, -- 0=unused, 1=used, 2=locked
    trade_no CHAR(22),
    note VARCHAR(32),
    mode INTEGER DEFAULT 0,
    category_id INTEGER DEFAULT 0,
    life INTEGER NOT NULL DEFAULT 1,
    use_life INTEGER NOT NULL DEFAULT 0,
    race VARCHAR(32)
);

CREATE INDEX idx_coupons_commodity_id ON coupons(commodity_id);
CREATE INDEX idx_coupons_owner ON coupons(owner);
CREATE INDEX idx_coupons_status ON coupons(status);

-- Managers table
CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    security_password VARCHAR(64),
    nickname VARCHAR(32),
    salt VARCHAR(32) NOT NULL,
    avatar VARCHAR(128),
    status INTEGER NOT NULL DEFAULT 0, -- 0=frozen, 1=normal
    type INTEGER NOT NULL DEFAULT 0, -- 0=system, 1=full_time, 2=day_time, 3=night_time
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_time DATETIME,
    last_login_time DATETIME,
    login_ip VARCHAR(128),
    last_login_ip VARCHAR(128),
    note VARCHAR(255)
);

-- Orders table
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner INTEGER NOT NULL DEFAULT 0, -- 0=guest, other=user_id
    user_id INTEGER NOT NULL DEFAULT 0, -- 0=system, other=merchant_id
    trade_no CHAR(19) NOT NULL UNIQUE,
    amount DECIMAL(10, 2) NOT NULL,
    commodity_id INTEGER NOT NULL,
    card_id INTEGER,
    card_num INTEGER NOT NULL DEFAULT 0,
    pay_id INTEGER NOT NULL,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_ip VARCHAR(64) NOT NULL,
    create_device INTEGER NOT NULL, -- 0=pc, 1=android, 2=ios, 3=ipad
    pay_time DATETIME,
    status INTEGER NOT NULL DEFAULT 0, -- 0=unpaid, 1=paid
    secret TEXT,
    password VARCHAR(32),
    contact VARCHAR(32),
    delivery_status INTEGER NOT NULL DEFAULT 0, -- 0=undelivered, 1=delivered
    pay_url VARCHAR(1024),
    coupon_id INTEGER,
    cost DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    from_user_id INTEGER, -- referrer user_id
    premium DECIMAL(10, 2) DEFAULT 0.00,
    widget TEXT,
    rent DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- cost price
    race VARCHAR(32),
    rebate DECIMAL(10, 2) DEFAULT 0.00,
    pay_cost DECIMAL(10, 2) DEFAULT 0.00,
    request_no CHAR(19) UNIQUE
);

CREATE INDEX idx_orders_commodity_id ON orders(commodity_id);
CREATE INDEX idx_orders_pay_id ON orders(pay_id);
CREATE INDEX idx_orders_contact ON orders(contact);
CREATE INDEX idx_orders_create_ip ON orders(create_ip);
CREATE INDEX idx_orders_owner ON orders(owner);
CREATE INDEX idx_orders_from_user_id ON orders(from_user_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_card_id ON orders(card_id);
CREATE INDEX idx_orders_create_time ON orders(create_time);
CREATE INDEX idx_orders_delivery_status ON orders(delivery_status);
CREATE INDEX idx_orders_coupon_id ON orders(coupon_id);

-- Order options table
CREATE TABLE order_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL UNIQUE,
    option TEXT,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Payment methods table
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(16) NOT NULL,
    icon VARCHAR(255),
    code VARCHAR(32) NOT NULL,
    commodity INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled for products
    recharge INTEGER NOT NULL DEFAULT 0, -- 0=disabled, 1=enabled for recharge
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    handle VARCHAR(64) NOT NULL,
    sort INTEGER NOT NULL DEFAULT 0,
    equipment INTEGER NOT NULL DEFAULT 0, -- 0=all, 1=mobile, 2=pc
    cost DECIMAL(10, 3) DEFAULT 0.000,
    cost_type INTEGER DEFAULT 0 -- 0=fixed, 1=percentage
);

CREATE INDEX idx_payments_commodity ON payments(commodity);
CREATE INDEX idx_payments_recharge ON payments(recharge);
CREATE INDEX idx_payments_sort ON payments(sort);
CREATE INDEX idx_payments_equipment ON payments(equipment);

-- Insert default payment methods
INSERT INTO payments (id, name, icon, code, commodity, recharge, handle, sort, equipment, cost, cost_type) VALUES
(1, '余额', '/assets/static/images/wallet.png', '#system', 1, 0, '#system', 999, 0, 0.000, 0),
(2, '支付宝', '/assets/user/images/cash/alipay.png', 'alipay', 1, 1, 'Epay', 1, 0, 0.000, 0);

-- Shared stores table
CREATE TABLE shared_stores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type INTEGER NOT NULL DEFAULT 0, -- 0=internal, others for extension
    name VARCHAR(128) NOT NULL,
    domain VARCHAR(128) NOT NULL UNIQUE,
    app_id VARCHAR(32) NOT NULL,
    app_key VARCHAR(64) NOT NULL,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    balance DECIMAL(14, 2) NOT NULL DEFAULT 0.00 -- cached balance
);

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(128) UNIQUE,
    phone VARCHAR(16) UNIQUE,
    qq VARCHAR(16),
    password VARCHAR(64) NOT NULL,
    salt VARCHAR(32) NOT NULL,
    app_key VARCHAR(32) NOT NULL,
    avatar VARCHAR(255),
    balance DECIMAL(14, 2) NOT NULL DEFAULT 0.00,
    coin DECIMAL(14, 2) NOT NULL DEFAULT 0.00, -- withdrawable coin
    integral INTEGER NOT NULL DEFAULT 0,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_time DATETIME,
    last_login_time DATETIME,
    login_ip VARCHAR(128),
    last_login_ip VARCHAR(128),
    pid INTEGER DEFAULT 0, -- parent user id
    recharge DECIMAL(14, 2) NOT NULL DEFAULT 0.00, -- total recharge
    total_coin DECIMAL(14, 2) NOT NULL DEFAULT 0.00, -- total earned coin
    status INTEGER NOT NULL DEFAULT 0, -- 0=banned, 1=normal
    business_level INTEGER, -- business level id
    nicename VARCHAR(10),
    alipay VARCHAR(64),
    wechat VARCHAR(255),
    settlement INTEGER NOT NULL DEFAULT 0 -- 0=alipay, 1=wechat
);

CREATE INDEX idx_users_pid ON users(pid);
CREATE INDEX idx_users_business_level ON users(business_level);
CREATE INDEX idx_users_coin ON users(coin);

-- User categories table
CREATE TABLE user_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    name VARCHAR(255), -- custom category name
    status INTEGER NOT NULL DEFAULT 0, -- 0=hidden, 1=visible
    UNIQUE(user_id, category_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_categories_status ON user_categories(status);

-- User commodities table
CREATE TABLE user_commodities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    commodity_id INTEGER NOT NULL,
    premium DECIMAL(10, 2) DEFAULT 0.00, -- markup
    name VARCHAR(255), -- custom name
    status INTEGER NOT NULL DEFAULT 0, -- 0=hidden, 1=visible
    UNIQUE(user_id, commodity_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (commodity_id) REFERENCES commodities(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_commodities_status ON user_commodities(status);

-- User groups table
CREATE TABLE user_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32) NOT NULL,
    icon VARCHAR(128),
    discount DECIMAL(4, 2) NOT NULL, -- discount percentage
    cost DECIMAL(4, 2) NOT NULL DEFAULT 0.00, -- commission percentage
    recharge DECIMAL(14, 2) NOT NULL UNIQUE -- required total recharge
);

-- Insert default user groups
INSERT INTO user_groups (id, name, icon, discount, cost, recharge) VALUES
(1, '一贫如洗', '/assets/static/images/group/ic_user level_1.png', 0.00, 0.30, 0.00),
(2, '小康之家', '/assets/static/images/group/ic_user level_2.png', 0.10, 0.25, 50.00),
(3, '腰缠万贯', '/assets/static/images/group/ic_user level_3.png', 0.20, 0.20, 100.00),
(4, '富甲一方', '/assets/static/images/group/ic_user level_4.png', 0.30, 0.15, 200.00),
(5, '富可敌国', '/assets/static/images/group/ic_user level_5.png', 0.40, 0.10, 300.00),
(6, '至尊', '/assets/static/images/group/ic_user level_6.png', 0.50, 0.05, 500.00);

-- User recharges table
CREATE TABLE user_recharges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trade_no CHAR(22) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    pay_id INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 0, -- 0=unpaid, 1=paid
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_ip VARCHAR(64) NOT NULL,
    pay_url VARCHAR(255),
    pay_time DATETIME,
    option TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_recharges_user_id ON user_recharges(user_id);
CREATE INDEX idx_user_recharges_pay_id ON user_recharges(pay_id);
CREATE INDEX idx_user_recharges_status ON user_recharges(status);

-- Manager logs table
CREATE TABLE manager_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(64) NOT NULL,
    nickname VARCHAR(32) NOT NULL,
    content VARCHAR(128) NOT NULL,
    create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_ip VARCHAR(64) NOT NULL,
    ua VARCHAR(255),
    risk INTEGER NOT NULL DEFAULT 0 -- 0=normal, 1=abnormal
);

CREATE INDEX idx_manager_logs_create_ip ON manager_logs(create_ip);
CREATE INDEX idx_manager_logs_create_time ON manager_logs(create_time);
CREATE INDEX idx_manager_logs_risk ON manager_logs(risk);
CREATE INDEX idx_manager_logs_email ON manager_logs(email);
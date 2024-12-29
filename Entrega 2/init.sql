-- Tabla usuario
CREATE TABLE usuario (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla sensor
CREATE TABLE sensor (
    id SERIAL PRIMARY KEY,
    tipo_sensor VARCHAR(100) NOT NULL,
    modelo VARCHAR(100) NOT NULL,
    ubicacion VARCHAR(255) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    fecha_instalacion TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla campo
CREATE TABLE campo (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    superficie FLOAT NOT NULL,
    tipo_cultivo VARCHAR(100) NOT NULL,
    ubicacion VARCHAR(255) NOT NULL,
    sensor_id INT REFERENCES sensor(id) ON DELETE SET NULL
);

-- Tabla tipo_plaga
CREATE TABLE tipo_plaga (
    tipo_plaga_id SERIAL PRIMARY KEY,
    nombre_comun VARCHAR(255) NOT NULL,
    descripcion TEXT,
    nombre_cientifico VARCHAR(255) NOT NULL
);

-- Tabla evento_plaga
CREATE TABLE evento_plaga (
    id SERIAL PRIMARY KEY,
    fecha_deteccion TIMESTAMP NOT NULL,
    ubicacion VARCHAR(255) NOT NULL,
    nivel_severidad VARCHAR(50) NOT NULL,
    acciones_tomadas TEXT,
    tipo_plaga_id INT REFERENCES tipo_plaga(tipo_plaga_id) ON DELETE CASCADE,
    campo_id INT REFERENCES campo(id) ON DELETE CASCADE,
    registro_vuelo_id INT -- Relación futura con registro_vuelo
);

-- Tabla medicion
CREATE TABLE medicion (
    id SERIAL PRIMARY KEY,
    fecha_hora TIMESTAMP NOT NULL,
    temperatura FLOAT NOT NULL,
    humedad FLOAT NOT NULL,
    luminosidad FLOAT NOT NULL,
    sensor_id INT REFERENCES sensor(id) ON DELETE CASCADE
);

-- Tabla configuracion
CREATE TABLE configuracion (
    id SERIAL PRIMARY KEY,
    umbra_temp FLOAT NOT NULL,
    umbra_humedad FLOAT NOT NULL,
    umbra_luminosidad FLOAT NOT NULL,
    usuario_id INT REFERENCES usuario(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla alert
CREATE TABLE alert (
    id SERIAL PRIMARY KEY,
    estado VARCHAR(50) NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    tipo_alerta VARCHAR(50) NOT NULL,
    usuario_id INT REFERENCES usuario(id) ON DELETE CASCADE,
    evento_plaga_id INT REFERENCES evento_plaga(id) ON DELETE CASCADE
);

-- Insertar usuarios
INSERT INTO usuario (nombre, email, password, rol, created_at, updated_at) VALUES
('Juan Pérez', 'juan@example.com', 'password123', 'admin', NOW(), NOW()),
('María García', 'maria@example.com', 'password456', 'usuario', NOW(), NOW()),
('Carlos López', 'carlos@example.com', 'password789', 'usuario', NOW(), NOW()),
('Ana Martínez', 'ana@example.com', 'passwordabc', 'usuario', NOW(), NOW()),
('Roberto Silva', 'roberto@example.com', 'passworddef', 'admin', NOW(), NOW()),
('Laura Torres', 'laura@example.com', 'passwordxyz', 'usuario', NOW(), NOW());

-- Insertar sensores
INSERT INTO sensor (tipo_sensor, modelo, ubicacion, estado, fecha_instalacion, created_at, updated_at) VALUES
('Temperatura', 'TempSensor2000', 'Sector Norte', 'Activo', '2024-01-01', NOW(), NOW()),
('Humedad', 'HumidPro', 'Sector Sur', 'Activo', '2024-01-02', NOW(), NOW()),
('Luminosidad', 'LuxMeter500', 'Sector Este', 'Activo', '2024-01-03', NOW(), NOW()),
('Temperatura', 'TempSensor3000', 'Sector Oeste', 'Inactivo', '2024-01-04', NOW(), NOW()),
('Humedad', 'HumidSense', 'Sector Central', 'Activo', '2024-01-05', NOW(), NOW()),
('Luminosidad', 'LuxMeter600', 'Sector Noroeste', 'Activo', '2024-01-06', NOW(), NOW());

-- Insertar campos
INSERT INTO campo (nombre, superficie, tipo_cultivo, ubicacion, sensor_id) VALUES
('Campo Maíz Norte', 150.5, 'Maíz', 'Norte', 1),
('Campo Soja Sur', 200.3, 'Soja', 'Sur', 2),
('Campo Trigo Este', 175.8, 'Trigo', 'Este', 3),
('Campo Girasol Oeste', 160.2, 'Girasol', 'Oeste', 4),
('Campo Maíz Central', 190.7, 'Maíz', 'Central', 5),
('Campo Soja Noroeste', 185.4, 'Soja', 'Noroeste', 6);

-- Insertar configuraciones
INSERT INTO configuracion (umbra_temp, umbra_humedad, umbra_luminosidad, usuario_id, created_at, updated_at) VALUES
(30.5, 75.0, 850.0, 1, NOW(), NOW()),
(28.5, 70.0, 800.0, 2, NOW(), NOW()),
(32.0, 80.0, 900.0, 3, NOW(), NOW()),
(29.0, 72.0, 825.0, 4, NOW(), NOW()),
(31.0, 77.0, 875.0, 5, NOW(), NOW()),
(27.5, 68.0, 775.0, 6, NOW(), NOW());

-- Insertar tipos de plagas
INSERT INTO tipo_plaga (nombre_comun, descripcion, nombre_cientifico) VALUES
('Gusano cogollero', 'Plaga común en cultivos de maíz', 'Spodoptera frugiperda'),
('Chinche verde', 'Afecta principalmente a la soja', 'Nezara viridula'),
('Pulgón del trigo', 'Plaga específica del trigo', 'Schizaphis graminum'),
('Isoca medidora', 'Afecta varios cultivos', 'Rachiplusia nu'),
('Trips', 'Plaga de pequeño tamaño', 'Frankliniella occidentalis'),
('Oruga de la hoja', 'Daña el follaje', 'Anticarsia gemmatalis');

-- Insertar eventos de plaga
INSERT INTO evento_plaga (fecha_deteccion, ubicacion, nivel_severidad, acciones_tomadas, tipo_plaga_id, campo_id, registro_vuelo_id) VALUES
('2024-01-10', 'Norte', 'Alto', 'Aplicación de insecticida', 1, 1, NULL),
('2024-01-11', 'Sur', 'Medio', 'Monitoreo intensivo', 2, 2, NULL),
('2024-01-12', 'Este', 'Bajo', 'Control biológico', 3, 3, NULL),
('2024-01-13', 'Oeste', 'Alto', 'Fumigación aérea', 4, 4, NULL),
('2024-01-14', 'Central', 'Medio', 'Aplicación de biocontroladores', 5, 5, NULL),
('2024-01-15', 'Noroeste', 'Bajo', 'Monitoreo preventivo', 6, 6, NULL);

-- Insertar alertas
INSERT INTO alert (estado, descripcion, fecha_hora, tipo_alerta, usuario_id, evento_plaga_id) VALUES
('Activa', 'Detección de plaga en sector norte', '2024-01-10 08:00:00', 'Urgente', 1, 1),
('Resuelta', 'Nivel de humedad crítico', '2024-01-11 09:30:00', 'Advertencia', 2, 2),
('Pendiente', 'Temperatura fuera de rango', '2024-01-12 10:45:00', 'Información', 3, 3),
('Activa', 'Nuevo brote de plaga detectado', '2024-01-13 11:15:00', 'Urgente', 4, 4),
('Resuelta', 'Luminosidad bajo nivel óptimo', '2024-01-14 12:30:00', 'Advertencia', 5, 5),
('Pendiente', 'Posible riesgo de infestación', '2024-01-15 13:45:00', 'Información', 6, 6);

-- Insertar mediciones
-- Insertar mediciones
INSERT INTO medicion (fecha_hora, temperatura, humedad, luminosidad, sensor_id) VALUES
('2024-01-10 08:00:00', 25.5, 60.0, 700.0, 1),
('2024-01-11 09:00:00', 26.0, 65.0, 750.0, 2),
('2024-01-12 10:00:00', 27.5, 70.0, 800.0, 3),
('2024-01-13 11:00:00', 28.0, 75.0, 850.0, 4),
('2024-01-14 12:00:00', 29.5, 80.0, 900.0, 5),
('2024-01-15 13:00:00', 30.0, 85.0, 950.0, 6);

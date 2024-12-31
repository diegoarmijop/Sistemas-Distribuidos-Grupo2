-- Crear la tabla para "nodo"
CREATE TABLE nodo (
    id SERIAL PRIMARY KEY,
    estado VARCHAR(255) NOT NULL,
    ubicacion VARCHAR(255) NOT NULL
);

-- Crear la tabla para "sensor"
CREATE TABLE sensor (
    id SERIAL PRIMARY KEY,
    temperatura VARCHAR(255) NOT NULL,
    humedad VARCHAR(255) NOT NULL,
    insectos VARCHAR(255) NOT NULL,
    luz VARCHAR(255) NOT NULL
);

-- Crear la tabla para "ruta"
CREATE TABLE ruta (
    id SERIAL PRIMARY KEY,
    fecha_hora_inicio TIMESTAMP NOT NULL,
    fecha_hora_termino TIMESTAMP NOT NULL,
    flag_dron VARCHAR(255) NOT NULL
);

-- Crear la tabla para "dron"
CREATE TABLE dron (
    id SERIAL PRIMARY KEY,
    estado VARCHAR(255) NOT NULL,
    modelo VARCHAR(255) NOT NULL,
    ubicacion VARCHAR(255) NOT NULL,
    ruta_id INTEGER,
    FOREIGN KEY (ruta_id) REFERENCES ruta(id) ON DELETE SET NULL
);

-- Insertar datos de ejemplo en "nodo"
INSERT INTO nodo (estado, ubicacion) VALUES
('Activo', 'Ubicación Nodo 2');

-- Insertar datos de ejemplo en "sensor" (3 sensores)
INSERT INTO sensor (temperatura, humedad, insectos, luz) VALUES
('22°C', '60%', 'Sí', 'Baja'),
('24°C', '50%', 'No', 'Alta'),
('20°C', '55%', 'No', 'Moderada');

-- Insertar datos de ejemplo en "ruta"
INSERT INTO ruta (fecha_hora_inicio, fecha_hora_termino, flag_dron) VALUES
('2024-12-29 09:00:00', '2024-12-29 11:00:00', 'Dron1'),
('2024-12-29 12:00:00', '2024-12-29 14:00:00', 'Dron2');

-- Insertar datos de ejemplo en "dron" (2 drones)
INSERT INTO dron (estado, modelo, ubicacion, ruta_id) VALUES
('activo', 'DJI-003', 'Sur', NULL),
('activo', 'DJI-004', 'Norte', NULL);

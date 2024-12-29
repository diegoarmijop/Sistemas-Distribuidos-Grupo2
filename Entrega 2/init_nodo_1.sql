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
('Activo', 'Ubicación Nodo 1');

-- Insertar datos de ejemplo en "sensor" (2 sensores)
INSERT INTO sensor (temperatura, humedad, insectos, luz) VALUES
('25°C', '50%', 'No', 'Alta'),
('23°C', '55%', 'Sí', 'Moderada');

-- Insertar datos de ejemplo en "ruta"
INSERT INTO ruta (fecha_hora_inicio, fecha_hora_termino, flag_dron) VALUES
('2024-12-29 08:00:00', '2024-12-29 10:00:00', 'Dron1');

-- Insertar datos de ejemplo en "dron"
INSERT INTO dron (estado, modelo, ubicacion, ruta_id) VALUES
('En vuelo', 'Modelo A', 'Ubicación Dron 1', 1);

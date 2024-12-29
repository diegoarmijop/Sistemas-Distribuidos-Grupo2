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
('Inactivo', 'Ubicación Nodo 3');

-- Insertar datos de ejemplo en "sensor" (2 sensores)
INSERT INTO sensor (temperatura, humedad, insectos, luz) VALUES
('21°C', '65%', 'Sí', 'Moderada'),
('26°C', '45%', 'No', 'Alta');

-- Insertar datos de ejemplo en "ruta"
INSERT INTO ruta (fecha_hora_inicio, fecha_hora_termino, flag_dron) VALUES
('2024-12-29 10:00:00', '2024-12-29 12:00:00', 'Dron3');

-- Insertar datos de ejemplo en "dron"
INSERT INTO dron (estado, modelo, ubicacion, ruta_id) VALUES
('En vuelo', 'Modelo C', 'Ubicación Dron 3', 1);

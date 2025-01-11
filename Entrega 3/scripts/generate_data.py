import random
import requests
import time
from datetime import datetime

def generate_sensor_data():
    return {
        "temperatura": f"{random.randint(20, 70)}°C",  # Generar un entero entre 20 y 35
        "humedad": f"{random.randint(30, 80)}%",      # Generar un entero entre 30 y 80
        "insectos": str(random.randint(0, 100)),      # Generar un entero entre 0 y 100
        "luz": f"{random.randint(500, 1000)} lux"     # Generar un entero entre 500 y 1000
    }

def send_to_sensor(sensor_id, data):
    try:
        url = f"http://localhost:8081/api/sensor/publicar/{sensor_id}"       
        response = requests.post(url, json=data)
        print(f"Datos enviados al Sensor {sensor_id}: {response.status_code}")
        print(f"Datos: {data}")
    except Exception as e:
        print(f"Error enviando datos al Sensor {sensor_id}: {e}")

def simulate_sensors():
    # Asumiendo que tienes sensores del 1 al 3
    sensor_ids = [1, 2]
    
    while True:
        for sensor_id in sensor_ids:
            data = generate_sensor_data()
            send_to_sensor(sensor_id, data)
            time.sleep(2)  # Esperar 2 segundos entre envíos.
        print("\n--- Nueva ronda de datos ---\n")

if __name__ == "__main__":
    print("Iniciando simulación de sensores...")
    simulate_sensors()

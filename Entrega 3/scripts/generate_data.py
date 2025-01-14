import random
import requests
import time
from datetime import datetime

def generate_insect_level():
    level = random.randint(0, 100)
    if level < 30:
        return "bajo"
    elif level < 60:
        return "medio"
    elif level < 90:
        return "alto"
    else:
        return "abundancia peligrosa"

def generate_sensor_data():
    # Generar temperatura entre 0 y 50°C
    temperatura = random.randint(0, 50)
    # Generar humedad entre 10 y 90%
    humedad = random.randint(10, 90)
    # Generar nivel de luz entre 1 y 15 UV
    luz = random.randint(1, 15)
    # Generar nivel de insectos
    insectos = generate_insect_level()

    return {
        "temperatura": f"{temperatura}°C",
        "humedad": f"{humedad}%",
        "insectos": insectos,
        "luz": f"{luz}UV"
    }

def print_condition_evaluation(data):
    # Extraer valores numéricos
    temperatura = int(data["temperatura"].replace("°C", ""))
    humedad = int(data["humedad"].replace("%", ""))
    luz = int(data["luz"].replace("UV", ""))
    insectos = data["insectos"]

    print("\nEvaluación de condiciones:")
    conditions = []

    # Evaluar temperatura
    if temperatura > 30:
        conditions.append("Temperatura alta")
        print("- Temperatura crítica alta")
    elif temperatura < 15:
        conditions.append("Temperatura baja")
        print("- Temperatura crítica baja")
    else:
        print("- Temperatura normal")

    # Evaluar humedad
    if humedad < 20 or humedad > 70:
        conditions.append("Humedad baja-alta")
        print("- Humedad crítica")
    else:
        print("- Humedad normal")

    # Evaluar insectos
    if insectos in ["alto", "abundancia peligrosa"]:
        conditions.append("Nivel alto de insectos")
        print("- Nivel de insectos crítico")
    else:
        print("- Nivel de insectos normal")

    # Evaluar luz
    if luz > 11:
        conditions.append("Nivel de luz extremadamente alto")
        print("- Nivel de luz extremadamente alto")
    elif luz > 8:
        conditions.append("Nivel de luz alto")
        print("- Nivel de luz alto")
    else:
        print("- Nivel de luz normal")

    if conditions:
        print(f"\nCondiciones críticas detectadas: {'/'.join(conditions)}")
    else:
        print("\nNo se detectaron condiciones críticas")

    return conditions

def send_to_sensor(sensor_id, data):
    try:
        url = f"http://localhost:8081/api/sensor/publicar/{sensor_id}"       
        response = requests.post(url, json=data)
        print(f"\nDatos enviados al Sensor {sensor_id}: {response.status_code}")
        print(f"Datos: {data}")
        
        # Evaluar condiciones después de enviar
        conditions = print_condition_evaluation(data)
        if conditions:
            print(f"¡ALERTA! Se generará una alerta con tipos: {'/'.join(conditions)}")
        
    except Exception as e:
        print(f"Error enviando datos al Sensor {sensor_id}: {e}")

def simulate_sensors():
    sensor_ids = [1, 2]
    
    while True:
        for sensor_id in sensor_ids:
            data = generate_sensor_data()
            send_to_sensor(sensor_id, data)
            time.sleep(5)  # Aumentado a 5 segundos para mejor lectura
        print("\n" + "="*50 + "\nNueva ronda de datos\n" + "="*50 + "\n")

if __name__ == "__main__":
    print("Iniciando simulación de sensores...")
    print("Presiona Ctrl+C para detener")
    try:
        simulate_sensors()
    except KeyboardInterrupt:
        print("\nSimulación detenida por el usuario")
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/ServiceStatus.css';

const ServiceStatus = () => {
    const [services, setServices] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchServices = async () => {
            try {
                const response = await axios.get('http://localhost:8083/health');
                setServices(response.data.services);
                setLoading(false);
            } catch (err) {
                setError('Error al obtener el estado de los servicios');
                setLoading(false);
            }
        };

        fetchServices();
        // Actualizar cada 30 segundos
        const interval = setInterval(fetchServices, 30000);
        return () => clearInterval(interval);
    }, []);

    if (loading) return <div className="services-status-container"><p>Cargando estado de los servicios...</p></div>;

    return (
        <div className="services-status-container">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Estado de los Microservicios</h1>
            {error && <p className="error-message">{error}</p>}
            <div className="services-grid">
                {services.map((service, index) => (
                    <div key={index} className={`service-card ${service.status}`}>
                        <h3>{service.name}</h3>
                        <p><strong>Estado:</strong> {service.status}</p>
                        <p><strong>Puerto:</strong> {service.port}</p>
                        <p><strong>Contenedor:</strong> {service.container}</p>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ServiceStatus;
import React, { useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';
import '../assets/styles/Home.css';

function Home() {
    const [cursos, setCursos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useContext(UserContext);
    const navigate = useNavigate();
    const isAdmin = user && user.user_type && user.user_type.toLowerCase() === 'administrador';

    useEffect(() => {
        const fetchCursos = async () => {
            try {
                const response = await axios.get('http://localhost:8080/courses');
                setCursos(response.data);
                setError(null);
            } catch (error) {
                console.log('Error fetching courses:', error.message);
                // En lugar de mostrar error, simplemente dejamos cursos como array vacío
                setCursos([]);
                setError(null);
            } finally {
                setLoading(false);
            }
        };
        fetchCursos();
    }, []);

    if (loading) {
        return (
            <div className="home-container">
                <div className="loading-message">Cargando...</div>
            </div>
        );
    }

    return (
        <div className="home-container">
            <main className="main-content">
                <section className="welcome-section">
                    <h1>Bienvenido{user ? `, ${user.username}` : ''}</h1>
                    <p className="description">Explora y administra tus cursos con facilidad. Aquí puedes encontrar información detallada sobre todos los cursos disponibles y gestionar tus cursos activos.</p>
                    <div className="quick-access-grid">
                        <div className="quick-card" onClick={() => navigate('/my-courses')}>
                            <h3>Mis Cursos</h3>
                            <p>Accede rápidamente a tus cursos inscritos.</p>
                        </div>
                        <div className="quick-card" onClick={() => navigate('/search')}>
                            <h3>Buscar Cursos</h3>
                            <p>Encuentra cursos por nombre o descripción.</p>
                        </div>
                        {isAdmin && (
                            <div className="quick-card" onClick={() => navigate('/manage-courses')}>
                                <h3>Gestión de Cursos</h3>
                                <p>Administra y edita los cursos de la plataforma.</p>
                            </div>
                        )}
                        {isAdmin && (
                            <div className="quick-card" onClick={() => navigate('/services-status')}>
                                <h3>Estado de Servicios</h3>
                                <p>Verifica el estado de los microservicios.</p>
                            </div>
                        )}
                    </div>
                </section>
                <section className="courses-section">
                    <h2>Cursos Disponibles</h2>
                    {Array.isArray(cursos) && cursos.length === 0 ? (
                        <div className="no-courses-message">
                            <p>No hay cursos disponibles.</p>
                            {isAdmin && (
                                <button onClick={() => navigate('/add-course')} className="add-course-button">
                                    Agregar un curso
                                </button>
                            )}
                        </div>
                    ) : (
                        Array.isArray(cursos) && cursos.length > 0 && (
                            <div className="courses-grid">
                                {cursos.map(curso => (
                                    <div key={curso.id} className="course-card">
                                        <div className="course-card-content">
                                            <h3>{curso.name}</h3>
                                            <p className="course-desc">{curso.description}</p>
                                        </div>
                                        <button className="details-button" onClick={() => navigate(`/courses/${curso.id}`)}>
                                            Ver detalles
                                        </button>
                                    </div>
                                ))}
                            </div>
                        )
                    )}
                </section>
            </main>
        </div>
    );
}

export default Home;

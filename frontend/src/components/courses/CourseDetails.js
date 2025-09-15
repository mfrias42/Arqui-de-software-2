import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate, Link, useParams } from 'react-router-dom';
import '../assets/styles/CourseDetails.css';


function CourseDetails() {
    const { courseId } = useParams();
    const navigate = useNavigate();
    const [course, setCourse] = useState(null);
    const [loading, setLoading] 
    = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchCourse = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/courses/${courseId}`);
                setCourse(response.data);
                setLoading(false);
            } catch (err) {
                setError('Error fetching course details: ' + err.message);
                setLoading(false);
            }
        };

        fetchCourse();
    }, [courseId]);

    const handleEnroll = async () => {
        const userId = localStorage.getItem('userId');
        if (!userId) {
            alert('Por favor, inicia sesión para inscribirte en el curso.');
            return;
        }

        try {
            // Intentar inscribirse en el curso
            await axios.post(`http://localhost:8085/inscriptions`, {
                user_id: parseInt(userId),
                course_id: parseInt(courseId)
            });
            alert('Inscripción exitosa!');
            navigate('/my-courses'); // Redirigir a "Mis Cursos" después de inscribirse
        } catch (err) {
            alert('Error en la inscripción: ' + err.response?.data?.error || err.message);
        }
    };

    if (loading) return <div className="course-details-container">Cargando...</div>;
    if (error) return <div className="course-details-container">{error}</div>;

    return (
        <div className="course-details-outer">
            <header className="course-details-header">
                <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            </header>
            <main className="course-details-main">
                <div className="course-details-card">
                    <div className="course-info-block">
                        {course.imageBase64 && (
                            <img src={course.imageBase64} alt="Imagen del curso" style={{maxWidth: 220, margin: '0 auto 18px auto', borderRadius: 10, display: 'block'}} />
                        )}
                        <h1>{course.name}</h1>
                        <p><strong>Descripción:</strong> {course.description}</p>
                        <p><strong>Categoría:</strong> {course.category}</p>
                        <p><strong>Duración:</strong> {course.duration}</p>
                        <p><strong>Instructor ID:</strong> {course.instructor_id}</p>
                        <p><strong>Capacidad:</strong> {course.capacity}</p>
                        <p><strong>Rating:</strong> {course.rating}</p>
                    </div>
                    <div className="course-actions-block">
                        {course.capacity > 0 ? (
                            <button onClick={handleEnroll} className="enroll-button">Inscribirse</button>
                        ) : (
                            <p className="full-message">El curso está lleno. No se puede inscribir.</p>
                        )}
                        <Link to={`/courses/${courseId}/files`} className="files-button">Ver Archivos</Link>
                        <Link to={`/courses/${courseId}/comments`} className="comments-button">Ver Comentarios</Link>
                    </div>
                </div>
            </main>
        </div>
    );
}

export default CourseDetails;

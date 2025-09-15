import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate } from 'react-router-dom';
import '../assets/styles/MyCourses.css';

function MyCourses() {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchMyCourses = async () => {
            try {
                const userId = localStorage.getItem('userId');
                if (!userId) {
                    throw new Error('User ID not found');
                }
                const inscriptionsResponse = await axios.get(`http://localhost:8085/users/${userId}/inscriptions`);
                const courseIds = inscriptionsResponse.data
                    .filter(inscription => inscription.course_id)
                    .map(inscription => inscription.course_id);
                if (courseIds.length === 0) {
                    setCourses([]);
                    return;
                }
            const coursesData = await Promise.all(
                courseIds.map(async (id) => {
                    try {
                        const response = await axios.get(`http://localhost:8080/courses/${id}`);
                        return response.data;
                    } catch (error) {
                        return null;
                    }
                })
            );
            setCourses(coursesData.filter(course => course !== null));
        } catch (err) {
            setError('Error fetching my courses');
        } finally {
            setLoading(false);
        }
    };
    fetchMyCourses();
}, []);

    const handleUploadClick = (courseId) => {
        navigate(`/upload/${courseId}`);
    };
    const handleCommentClick = (courseId) => {
        navigate(`/courses/${courseId}/comments`);
    };
    if (loading) return <div className="my-courses-container">Cargando...</div>;
    if (error) return <div className="my-courses-container">{error}</div>;
    return (
        <div className="my-courses-outer">
            <header className="my-courses-header">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Mis Cursos</h1>
            </header>
            <main className="my-courses-main">
            {courses.length > 0 ? (
                    <div className="my-courses-grid">
                        {courses.map(course => (
                            <div key={course.id} className="my-course-card">
                                <div className="my-course-card-content">
                                    <h3>{course.name}</h3>
                                    <p className="my-course-desc">{course.description}</p>
                                </div>
                                <div className="my-course-actions">
                                    <Link to={`/courses/${course.id}`} className="details-button">Ver detalles</Link>
                                <button onClick={() => handleUploadClick(course.id)} className="upload-button">Subir Archivo</button>
                                    <button onClick={() => handleCommentClick(course.id)} className="comment-button">Comentar</button>
                                </div>
                            </div>
                        ))}
                    </div>
            ) : (
                <p>No estás inscrito en ningún curso aún.</p>
            )}
            </main>
        </div>
    );
}

export default MyCourses;

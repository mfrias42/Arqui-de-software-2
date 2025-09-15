import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import CommentForm from './CommentForm';
import CommentList from './CommentList';
import '../assets/styles/CourseComments.css';

const CourseComments = () => {
    const { courseId } = useParams();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        // Aquí podrías hacer una llamada para obtener comentarios si es necesario
        setLoading(false); // Simulación de carga
    }, [courseId]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div className="course-comments-outer">
            <div className="course-comments-card">
            <CommentForm />
            <CommentList />
            </div>
        </div>
    );
};

export default CourseComments;

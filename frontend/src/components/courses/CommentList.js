import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import '../assets/styles/CommentList.css';

function CommentList() {
    const { courseId } = useParams();
    const [comments, setComments] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchComments = async () => {
            const token = localStorage.getItem('token');
            try {
                const response = await axios.get(
                    `http://localhost:8080/courses/${courseId}/comments`,
                    {
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json'
                        }
                    }
                );
                setComments(response.data);
            } catch (err) {
                console.error('Error al obtener los comentarios:', err);
                setError(err.response?.data?.error || err.message);
            }
        };

        fetchComments();
    }, [courseId]);

    return (
        <div className="comment-list-container">
            <h1>Comentarios del Curso</h1>
            {error && <p className="error-message">{error}</p>}
            {comments.length === 0 ? (
                <p className="no-comments-message">Aún no hay comentarios para este curso.</p>
            ) : (
                <ul>
                    {comments.map((comment) => (
                        <li key={comment.id} className="comment-item">
                            <p>{comment.content}</p>
                            <small>Usuario: {comment.user_id}</small>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default CommentList;

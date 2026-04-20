import React from 'react';

const PublicLayout: React.FC<{ children: React.ReactNode }> = ({children}) => {
    console.log('>>> PublicLayout render');
    return (
        <div style={{ minHeight: '100dvh', padding: '1rem' }}>
            {children}
        </div>
    );
};

export default PublicLayout;

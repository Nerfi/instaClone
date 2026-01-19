graph TD
Start([Usuario solicita reset]) --> FP1[POST /forgot-password]

    FP1 --> FP2[AuthHandlers.ForgotPassword]
    FP2 --> FP3[AuthService.FindUserByEmail]
    
    FP3 --> FP4{¿Usuario existe?}
    FP4 -->|No| FP5[Respuesta genérica]
    FP4 -->|Sí| FP6[token.GenerateResetToken]
    
    FP6 --> FP7[AuthService.SavePasswordResetToken]
    FP7 --> FP8[AuthRepository.SavePasswordResetToken]
    FP8 --> FP9[UPDATE password_resets SET used=TRUE]
    FP9 --> FP10[INSERT INTO password_resets]
    
    FP10 --> FP11[EmailService.SendPasswordResetEmail]
    FP11 --> FP5
    FP5 --> End1([Fin: Token enviado por email])
    
    End1 -.Usuario recibe email.-> Start2([Usuario hace click en link])
    
    Start2 --> RP1[POST /reset-password<br/>token + new_password]
    RP1 --> RP2[AuthHandlers.ResetPassword]
    
    RP2 --> RP3[token.HashToken]
    RP3 --> RP4[AuthService.ValidateResetToken]
    RP4 --> RP5[AuthRepository.ValidateResetToken]
    
    RP5 --> RP6{¿Token válido,<br/>no usado,<br/>no expirado?}
    RP6 -->|No| RP7[Error: Invalid or expired token]
    RP6 -->|Sí| RP8[security.HashPassword]
    
    RP8 --> RP9[AuthService.UpdateUserPassword]
    RP9 --> RP10[AuthRepository.UpdateUserPassword]
    RP10 --> RP11[UPDATE users SET password]
    
    RP11 --> RP12[AuthService.MarkTokenAsUsed]
    RP12 --> RP13[AuthRepository.MarkTokenAsUsed]
    RP13 --> RP14[UPDATE password_resets SET used=TRUE]
    
    RP14 --> RP15[Success: Password reset]
    RP7 --> End2([Fin: Error])
    RP15 --> End3([Fin: Contraseña actualizada])
    
    style FP2 fill:#e1f5ff
    style RP2 fill:#fff4e1
    style FP8 fill:#e8f5e9
    style RP10 fill:#e8f5e9
    style RP13 fill:#e8f5e9

flowchart TD

    %% --- FORGOT PASSWORD FLOW ---
    Start([Usuario solicita reset]) --> FP1[POST /forgot-password]
    FP1 --> FP2[AuthHandlers.ForgotPassword]
    FP2 --> FP3[AuthService.FindUserByEmail]

    FP3 -->|No existe| FP5[Respuesta genérica]
    FP3 -->|Sí existe| FP6[Generar token de reset]

    FP6 --> FP7[Guardar token]
    FP7 --> FP8[Marcar tokens previos como usados]
    FP8 --> FP9[Insertar nuevo registro en password_resets]

    FP9 --> FP10[Enviar email con link de reset]
    FP10 --> FP5
    FP5 --> End1([Fin: Token enviado por email])

    End1 -. Usuario recibe email .-> Start2([Usuario hace click en link])

    %% --- RESET PASSWORD FLOW ---
    Start2 --> RP1[POST /reset-password (token + nueva contraseña)]
    RP1 --> RP2[AuthHandlers.ResetPassword]
    RP2 --> RP3[Hashear token recibido]
    RP3 --> RP4[Validar token en BD]

    RP4 -->|Token inválido / usado / expirado| RP7[Error: Invalid or expired token]
    RP4 -->|Token válido| RP8[Hashear nueva contraseña]

    RP8 --> RP9[Actualizar contraseña del usuario]
    RP9 --> RP10[Marcar token como usado]
    RP10 --> RP11([Fin: Contraseña actualizada])

    RP7 --> EndError([Fin: Error])

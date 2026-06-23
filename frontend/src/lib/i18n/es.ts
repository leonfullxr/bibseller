// Spanish (Castilian) catalogue - M8.2 (#46), D4 Spain-first. Typed as a full
// Record<MessageKey, string>, so the compiler enforces that every en key is
// present and rejects extras (the key-set check is also asserted at runtime in
// locale.test.ts). Informal "tú"; a race bib is a "dorsal". `{...}` placeholders
// match en. Language names (lang.en/lang.es) stay as their own endonyms.
import type { MessageKey } from './en';

export const es: Record<MessageKey, string> = {
	'lang.en': 'English',
	'lang.es': 'Español',
	'lang.switch': 'Idioma',

	'nav.races': 'Carreras',
	'nav.sell': 'Vender',
	'nav.myListings': 'Mis anuncios',
	'nav.inbox': 'Mensajes',
	'nav.login': 'Iniciar sesión',
	'nav.register': 'Registrarse',
	'nav.logout': 'Cerrar sesión',
	'banner.verifyEmail': 'Verifica tu correo para poder vender y chatear.',
	'banner.resend': 'Reenviar correo',
	'banner.verifySent': 'Correo de verificación enviado: revisa tu bandeja de entrada (y spam).',
	'banner.suggestText': 'Esta página está disponible en español.',
	'banner.suggestAccept': 'Ver en español',
	'banner.suggestDismiss': 'Ahora no',
	'footer.tagline': 'Sin comisiones, en toda la UE. Sin ánimo de lucro por diseño.',
	'footer.github': 'GitHub',
	'footer.terms': 'Términos',
	'footer.privacy': 'Privacidad',
	'footer.contact': 'Contacto',

	'home.title': 'Bibseller - los dorsales encuentran nuevos corredores',
	'home.metaDescription':
		'Mercado sin ánimo de lucro, en toda la UE, que conecta a corredores que no pueden participar con quienes se quedaron sin inscripción.',
	'home.heroTitle': 'Los dorsales encuentran',
	'home.heroTitleHighlight': 'nuevos corredores',
	'home.tagline':
		'¿Lesionado? ¿Te cambiaron los planes? ¿Te quedaste sin inscripción? Un mercado sin comisiones que conecta a quienes venden y compran dorsales, siempre según las normas de cada carrera.',
	'home.searchPlaceholder': 'Busca una carrera o ciudad…',
	'home.search': 'Buscar',
	'home.browseAll': 'o explora todas las carreras',
	'home.apiUnreachable': 'API no disponible - ejecuta',
	'home.upcoming': 'Próximas carreras',
	'home.seeAll': 'Ver todas',
	'home.howTitle': 'Cómo funciona',
	'home.step1Title': 'Publica tu dorsal',
	'home.step1Desc':
		'Publica tu inscripción en minutos. La política de transferencia se ve desde el principio.',
	'home.step2Title': 'Escribe al vendedor',
	'home.step2Desc':
		'Chatea de forma segura en Bibseller para acordar los detalles. Sin comisiones, siempre.',
	'home.step3Title': 'Organiza la transferencia',
	'home.step3Desc':
		'Reinscríbete a través del organizador cuando haga falta, o entrega el dorsal directamente.',
	'home.step4Title': 'Confirma la entrega',
	'home.step4Desc': 'Lo confirmáis los dos en el chat. El dorsal es suyo y habéis terminado.',
	'home.howNote': 'Algunas carreras son solo chat o reventa oficial, según el organizador.',
	'home.journeyTitle': 'El recorrido de comprador y vendedor',
	'home.journeyLead': 'Seis pasos sencillos, desde publicar un dorsal hasta la salida.',
	'home.journeySeller': 'Vendedor',
	'home.journeyBuyer': 'Comprador',
	'home.j1Title': 'Publica el dorsal',
	'home.j2Title': 'Lo encuentra y escribe',
	'home.j3Title': 'Responde y acuerda',
	'home.j4Title': 'Organiza la transferencia',
	'home.j5Title': 'Lo entrega',
	'home.j6Title': 'A la línea de salida',
	'home.contactTitle': 'Ponte en contacto',
	'home.contactLead':
		'¿Una carrera que añadir, una duda o sugerencias? Nos encantaría saber de ti.',
	'home.contactCta': 'Escríbenos',
	'home.modePlatformSaleName': 'Venta en la plataforma',
	'home.modePlatformSaleDesc':
		'La carrera permite la reventa: publica, chatea y paga de forma segura a través de la plataforma.',
	'home.modeOfficialName': 'Proceso oficial',
	'home.modeOfficialDesc':
		'La carrera gestiona su propio cambio de titular: te ponemos en contacto y enlazamos el trámite oficial.',
	'home.modeConnectName': 'Solo contacto',
	'home.modeConnectDesc':
		'Carreras restringidas o sin verificar: ponemos el chat, el resto queda entre vosotros.',
	'home.underConstruction': 'En construcción - sigue la',
	'home.roadmap': 'hoja de ruta',

	'policy.label.platform_sale': 'Reventa permitida',
	'policy.label.official_only': 'Cambio oficial',
	'policy.label.connect_only': 'Solo chat',
	'policy.label.unknown': 'Política sin verificar',
	'policy.disclaimer.platform_sale.title': 'Esta carrera permite la reventa de dorsales.',
	'policy.disclaimer.platform_sale.body':
		'Acuerda los detalles con quien vende en el chat y paga de forma segura a través de la plataforma: el dinero queda retenido hasta confirmar la transferencia. Sin comisiones.',
	'policy.disclaimer.official_only.title':
		'Esta carrera gestiona su propio proceso oficial de cambio de titular.',
	'policy.disclaimer.official_only.body':
		'Encontraos y acordad los detalles aquí; la transferencia en sí (y cualquier tasa oficial) se hace a través de la organización de la carrera. La plataforma nunca gestiona dinero para esta carrera.',
	'policy.disclaimer.connect_only.title': 'Esta carrera restringe la transferencia de dorsales.',
	'policy.disclaimer.connect_only.body':
		'La plataforma solo os pone en contacto: aquí no gestiona dinero ni se responsabiliza de ningún acuerdo entre tú y la otra parte. Se aplican las normas de la propia carrera; consúltalas antes de acordar nada.',
	'policy.disclaimer.unknown.title':
		'Política de transferencia aún sin verificar; trata esta carrera como solo chat.',
	'policy.disclaimer.unknown.body':
		'La plataforma solo os pone en contacto: aquí no gestiona dinero ni se responsabiliza de ningún acuerdo entre tú y la otra parte. Se aplican las normas de la propia carrera; consúltalas antes de acordar nada.',
	'policy.officialLink': 'Proceso de cambio oficial',

	'sport.running': 'Running',
	'sport.trail': 'Trail',
	'sport.triathlon': 'Triatlón',
	'sport.cycling': 'Ciclismo',
	'sport.obstacle': 'Obstáculos',
	'sport.other': 'Otro',

	'races.title': 'Explorar carreras - Bibseller',
	'races.metaDescription': 'Encuentra dorsales a la venta en eventos de running por toda la UE.',
	'races.heading': 'Explorar carreras',
	'races.filter.search': 'Buscar',
	'races.filter.searchPlaceholder': 'Carrera o ciudad…',
	'races.filter.country': 'País',
	'races.filter.sport': 'Deporte',
	'races.filter.policy': 'Política de transferencia',
	'races.filter.all': 'Todas',
	'races.filter.submit': 'Filtrar',
	'races.empty': 'Ninguna carrera coincide con esos filtros.',
	'races.clearFilters': 'Borrar filtros',
	'races.nextPage': 'Página siguiente ->',
	'races.mapHeading': 'Carreras por Europa',
	'races.mapHint': 'Los países resaltados tienen carreras: toca uno para filtrar.',
	'races.mapCountry': '{country}: {n} carreras',

	'raceCard.bibs.one': '{n} dorsal publicado',
	'raceCard.bibs.other': '{n} dorsales publicados',
	'listingCard.priceOnRequest': 'Precio a consultar',
	'listingCard.belowFace': 'por debajo del valor nominal',
	'listingCard.listedBy': 'Publicado por {name}',

	'raceDetail.title': '{name} - dorsales a la venta - Bibseller',
	'raceDetail.metaDescription': 'Dorsales para {name} ({date}, {city}).',
	'raceDetail.back': 'Volver a todas las carreras',
	'raceDetail.website': 'Web de la carrera',
	'raceDetail.bibsForSale.one': '{n} dorsal a la venta',
	'raceDetail.bibsForSale.other': '{n} dorsales a la venta',
	'raceDetail.sellCta': 'Vende tu dorsal',
	'raceDetail.empty': 'Aún no hay dorsales publicados para esta carrera.',
	'raceDetail.emptyHint': '¿Vendes el tuyo? Pronto podrás publicarlo.',

	'listingDetail.title': 'Dorsal para {name} - Bibseller',
	'listingDetail.back': 'Volver a {name}',
	'listingDetail.heading': 'Dorsal para {name}',
	'listingDetail.unavailable': 'Este anuncio ya no está disponible ({status}).',
	'listingDetail.listedByOn': 'Publicado por {name} el {date}',
	'listingDetail.contact': 'Contactar con quien vende',
	'listingDetail.toMessageSeller': 'para enviar un mensaje a quien vende.',
	'listingDetail.verifyToMessage': 'Verifica tu correo para enviar un mensaje a quien vende.',
	'listingDetail.accountSettings': 'Ajustes de la cuenta',
	'listingDetail.ownPre': 'Este es tu anuncio; gestiónalo desde',
	'listingDetail.yourListings': 'tus anuncios',
	'listingDetail.messageAria': 'Mensaje para quien vende',
	'listingDetail.messagePlaceholder': 'Hola, ¿sigue disponible este dorsal?',
	'listingDetail.ackText':
		'Entiendo que la plataforma no gestiona dinero ni se responsabiliza de esta transferencia; se aplican las normas de la propia carrera.',
	'listingDetail.send': 'Enviar mensaje',

	'report.summary': 'Denunciar este anuncio',
	'report.reason.forbidden_transfer': 'Transferencia no permitida',
	'report.reason.scam': 'Estafa',
	'report.reason.offensive': 'Ofensivo',
	'report.reason.other': 'Otro',
	'report.reasonAria': 'Motivo de la denuncia',
	'report.detailsAria': 'Detalles de la denuncia (opcional)',
	'report.detailsPlaceholder': 'Detalles (opcional)',
	'report.success': 'Gracias, hemos recibido tu denuncia de este anuncio.',
	'report.failed': 'No se pudo enviar la denuncia. Inténtalo de nuevo.',
	'report.networkError': 'Error de red. Inténtalo de nuevo.',
	'report.submitting': 'Enviando...',
	'report.submit': 'Enviar denuncia',

	'listingCta.buy': 'Compra segura - próximamente',
	'listingCta.buyTitle': 'El pago seguro llega con los pagos (M6)',

	'inbox.title': 'Mensajes - Bibseller',
	'inbox.heading': 'Mensajes',
	'inbox.emptyPre': 'Aún no tienes conversaciones. Explora',
	'inbox.emptyRacesLink': 'carreras',
	'inbox.emptyPost': 'y contacta con quien vende para iniciar una.',
	'role.seller': 'vendedor',
	'role.buyer': 'comprador',

	'chat.title': 'Chat con {name} - Bibseller',
	'chat.back': 'Volver a mensajes',
	'chat.about': 'sobre',
	'chat.block': 'Bloquear',
	'chat.unblock': 'Desbloquear',
	'chat.sharedImage': 'Imagen compartida',
	'chat.reportMsg': 'denunciar',
	'chat.messageAria': 'Tu mensaje',
	'chat.messagePlaceholder': 'Escribe un mensaje o adjunta una imagen...',
	'chat.attachAria': 'Adjuntar una imagen (JPEG o PNG)',
	'chat.send': 'Enviar',
	'chat.sending': 'Enviando...',
	'chat.imageTooLarge': 'Esa imagen es demasiado grande (5 MB máx.).',
	'chat.tooFast': 'Estás enviando mensajes demasiado rápido; espera un momento.',
	'chat.sendFailed': 'No se pudo enviar tu mensaje. Inténtalo de nuevo.',
	'chat.networkError': 'Error de red; comprueba tu conexión.',
	'chat.blockConfirm': '¿Bloquear a {name}? Ninguno de los dos podréis enviaros mensajes.',
	'chat.blocked': 'Usuario bloqueado.',
	'chat.blockFailed': 'No se pudo bloquear al usuario.',
	'chat.unblocked': 'Usuario desbloqueado.',
	'chat.unblockFailed': 'No se pudo desbloquear al usuario.',
	'chat.networkRetry': 'Error de red. Inténtalo de nuevo.',
	'chat.reportConfirm': '¿Denunciar este mensaje a los moderadores?',
	'chat.messageReported': 'Mensaje denunciado.',
	'chat.messageReportFailed': 'No se pudo denunciar el mensaje.',

	'sell.title': 'Vender un dorsal - Bibseller',
	'sell.heading': 'Vender un dorsal',
	'sell.lede':
		'Encuentra tu carrera y publica tu dorsal. Tú pones el precio (con el valor nominal como tope).',
	'sell.verifyNotice':
		'Verifica tu correo para publicar un anuncio; mientras, puedes buscar tu carrera abajo.',
	'sell.searchAria': 'Buscar carreras por nombre o ciudad',
	'sell.emptyPre': 'Ninguna próxima carrera coincide. Prueba otra búsqueda o',
	'sell.browseAllLink': 'explora todas las carreras',
	'sell.sellHere': 'Vender aquí',

	'sellForm.title': 'Publica tu dorsal para {name} - Bibseller',
	'sellForm.back': 'Volver a la búsqueda de carreras',
	'sellForm.heading': 'Publica tu dorsal',
	'sellForm.verifyNotice': 'Verifica tu correo para publicar un anuncio.',
	'sellForm.publish': 'Publicar anuncio',

	'listingFields.price': 'Precio de venta (EUR)',
	'listingFields.pricePlaceholder': 'p. ej. 45',
	'listingFields.original': 'Precio original / valor nominal (EUR)',
	'listingFields.optional': 'opcional',
	'listingFields.hint':
		'Indica el valor nominal y tu precio de venta no podrá superarlo: sin sobreprecio.',
	'listingFields.description': 'Descripción',
	'listingFields.descriptionPlaceholder': 'opcional - talla, detalles de entrega, etc.',

	'myListings.title': 'Mis anuncios - Bibseller',
	'myListings.heading': 'Mis anuncios',
	'myListings.emptyPre': 'Aún no tienes anuncios.',
	'myListings.listABib': 'Publica un dorsal',
	'myListings.edit': 'Editar',
	'myListings.cancel': 'Cancelar',

	'editListing.title': 'Editar anuncio - Bibseller',
	'editListing.back': 'Volver a mis anuncios',
	'editListing.heading': 'Editar anuncio',
	'editListing.save': 'Guardar cambios',

	'auth.email': 'Correo electrónico',
	'auth.password': 'Contraseña',

	'login.title': 'Iniciar sesión - Bibseller',
	'login.forgot': '¿Olvidaste tu contraseña?',
	'login.newHere': '¿Nuevo por aquí?',
	'login.createAccount': 'Crea una cuenta',

	'register.title': 'Crear cuenta - Bibseller',
	'register.heading': 'Crear cuenta',
	'register.displayName': 'Nombre visible',
	'register.haveAccount': '¿Ya tienes cuenta?',

	'forgot.title': 'Restablecer contraseña - Bibseller',
	'forgot.heading': 'Restablece tu contraseña',
	'forgot.lede': 'Introduce tu correo y te enviaremos un enlace para restablecerla.',
	'forgot.sent':
		'Si existe una cuenta con esa dirección, te hemos enviado un enlace para restablecer la contraseña. Revisa tu bandeja de entrada.',
	'forgot.submit': 'Enviar enlace',
	'forgot.backToLogin': 'Volver a iniciar sesión',

	'reset.title': 'Establecer una nueva contraseña - Bibseller',
	'reset.heading': 'Establece una nueva contraseña',
	'reset.done':
		'Tu contraseña se ha actualizado. Se ha cerrado la sesión en todos los dispositivos; inicia sesión con tu nueva contraseña.',
	'reset.missingToken': 'A este enlace le falta el token. Solicita uno nuevo.',
	'reset.requestLink': 'Solicitar un enlace',
	'reset.newPassword': 'Nueva contraseña',
	'reset.confirmPassword': 'Confirmar contraseña',
	'reset.submit': 'Actualizar contraseña',

	'verify.title': 'Verificar correo - Bibseller',
	'verify.okHeading': 'Correo verificado',
	'verify.okBody': 'Tu dirección de correo está confirmada; todo listo.',
	'verify.continue': 'Continuar',
	'verify.invalidHeading': 'Enlace no válido o caducado',
	'verify.invalidBody':
		'Este enlace de verificación ya no es válido. Inicia sesión y solicita uno nuevo.',
	'verify.signIn': 'Iniciar sesión',
	'verify.missingHeading': 'Nada que verificar',
	'verify.missingBody': 'Abre el enlace de verificación de tu correo para confirmar tu dirección.',
	'verify.home': 'Inicio',
	'verify.errorHeading': 'Algo salió mal',
	'verify.errorBody':
		'No hemos podido verificar tu correo ahora mismo. Inténtalo de nuevo en un momento.',

	'settings.title': 'Ajustes - Bibseller',
	'settings.heading': 'Ajustes',
	'settings.profile': 'Perfil',
	'settings.country': 'País',
	'settings.countryNotSet': 'Sin definir',
	'settings.profileUpdated': 'Perfil actualizado.',
	'settings.save': 'Guardar',
	'settings.password': 'Contraseña',
	'settings.currentPassword': 'Contraseña actual',
	'settings.confirmNewPassword': 'Confirmar nueva contraseña',
	'settings.passwordChanged':
		'Contraseña cambiada. Se ha cerrado la sesión en los demás dispositivos.',
	'settings.changePassword': 'Cambiar contraseña',
	'settings.sessions': 'Sesiones',
	'settings.sessionsNote':
		'Cierra la sesión de Bibseller en todos los dispositivos, incluido este.',
	'settings.logoutAll': 'Cerrar sesión en todos los dispositivos',
	'settings.deleteAccount': 'Eliminar cuenta',
	'settings.deleteNote':
		'Elimina permanentemente tu cuenta y sus datos. Disponible cuando llegue el cumplimiento completo del RGPD (M7).',
	'settings.deleteTitle': 'La eliminación de cuenta llega con confianza y seguridad (M7)',
	'settings.deleteSoon': 'Eliminar cuenta - próximamente',

	'error.notFound': 'Esa página no existe.',
	'error.generic': 'Algo salió mal.',
	'error.backHome': 'Volver al inicio',

	// Errores de la API (#49): se eligen por el `code` estable del sobre, no por
	// su mensaje en inglés; un code sin entrada recae en apiError.unknown.
	'apiError.unknown': 'Algo salió mal. Inténtalo de nuevo.',
	'apiError.unreachable': 'No se puede contactar con la API.',
	'apiError.not_found': 'No encontrado.',
	'apiError.invalid_parameter': 'Solicitud no válida.',
	'apiError.email_taken': 'Ya existe una cuenta con este correo.',
	'apiError.race_past': 'Esta carrera ya se ha celebrado.',
	'apiError.loadFailed': 'No se pudo cargar esta página. Inténtalo de nuevo.',

	// Errores de formulario/validación escritos en el frontend (server actions).
	'formError.invalidEmail': 'Introduce un correo electrónico válido.',
	'formError.displayNameLength': 'El nombre visible debe tener entre 2 y 50 caracteres.',
	'formError.passwordTooShort': 'La contraseña debe tener al menos 8 caracteres.',
	'formError.newPasswordTooShort': 'La nueva contraseña debe tener al menos 8 caracteres.',
	'formError.passwordMismatch': 'Las dos contraseñas no coinciden.',
	'formError.newPasswordMismatch': 'Las dos contraseñas nuevas no coinciden.',
	'formError.loginRequired': 'Introduce tu correo y tu contraseña.',
	'formError.invalidCredentials': 'Correo o contraseña incorrectos.',
	'formError.loginFailed': 'No se pudo iniciar sesión.',
	'formError.currentPasswordWrong': 'Tu contraseña actual es incorrecta.',
	'formError.changePasswordFailed': 'No se pudo cambiar la contraseña.',
	'formError.resetTokenMissing': 'A este enlace de restablecimiento le falta el token.',
	'formError.resetTokenInvalid': 'Este enlace de restablecimiento no es válido o ha caducado.',
	'formError.resetFailed': 'No se pudo restablecer la contraseña.',
	'formError.resetEmailFailed':
		'No se pudo enviar el correo de restablecimiento. Inténtalo de nuevo.',
	'formError.tooManyRequests': 'Demasiadas solicitudes. Espera un minuto e inténtalo de nuevo.',
	'formError.missingRace': 'Falta la carrera; vuelve a empezar desde la página de la carrera.',
	'formError.missingListingId': 'Falta el identificador del anuncio.',
	'formError.invalidAmount': 'Introduce un importe válido, p. ej. 45 o 45.00.',
	'formError.priceExceedsFace': 'El precio de venta no puede superar el valor nominal original.',
	'formError.verifyToContact': 'Verifica tu correo antes de contactar con quien vende.',
	'formError.verifyToPublish': 'Verifica tu correo antes de publicar un anuncio.',
	'formError.emptyMessage': 'Escribe primero un mensaje para quien vende.',
	'formError.ackRequired': 'Acepta las condiciones para continuar.',
	'formError.ackFailed': 'No se pudo registrar tu aceptación. Inténtalo de nuevo.',
	'formError.cannotContact': 'No puedes contactar con quien vende para este anuncio.',
	'formError.listingUnavailable': 'Este anuncio ya no está disponible.',
	'formError.contactFailed': 'No se pudo iniciar la conversación.',
	'formError.editOwnOnly': 'Solo puedes editar tus propios anuncios.',
	'formError.editNotActive': 'Este anuncio ya no está activo y no se puede editar.',
	'formError.editFailed': 'No se pudo actualizar el anuncio.',
	'formError.cancelFailed': 'No se pudo cancelar el anuncio.',

	// Páginas legales (M7-lite, #11). Copia BORRADOR pendiente de revisión legal (D10).
	'legal.draftNotice': 'Borrador - pendiente de revisión legal y aún no vinculante.',
	'legal.terms.title': 'Términos del servicio - Bibseller',
	'legal.terms.metaDescription': 'Las condiciones que rigen el uso del mercado de Bibseller.',
	'legal.terms.heading': 'Términos del servicio',
	'legal.terms.intro':
		'Bibseller es un mercado sin comisiones que conecta a quienes quieren ceder un dorsal con quienes buscan uno. Ponemos el canal; no somos parte de ninguna transferencia entre usuarios.',
	'legal.terms.roleHeading': 'Nuestro papel',
	'legal.terms.roleBody':
		'Nunca cobramos comisión y solo repercutimos las tarifas del procesador de pagos cuando se ofrece el pago seguro. Lo que puedes hacer con un dorsal lo decide cada organización; la plataforma aplica esas normas en cada carrera.',
	'legal.terms.policyHeading': 'Responsabilidades según la política de transferencia de la carrera',
	'legal.terms.policyIntro':
		'Cada carrera tiene una política de transferencia que determina qué hace la plataforma y quién es responsable de qué:',
	'legal.terms.policy.platform_sale':
		'Reventa permitida: la carrera permite las transferencias. Acuerdas los detalles en el chat y, cuando llegue el pago seguro, podrás pagar a través de la plataforma; la transferencia del dorsal y la inscripción siguen sujetas al proceso de la organización.',
	'legal.terms.policy.official_only':
		'Proceso oficial: la carrera gestiona su propio cambio de titular. Te ponemos en contacto y enlazamos ese proceso; la transferencia y cualquier tasa oficial se realizan con la organización, nunca a través de la plataforma.',
	'legal.terms.policy.connect_only':
		'Solo contacto: la carrera restringe las transferencias. La plataforma solo os pone en contacto: no gestiona dinero ni participa en ningún acuerdo. Se aplican las normas de la carrera y confirmas que lo entiendes antes de contactar con quien vende.',
	'legal.terms.policy.unknown':
		'Sin verificar: hasta que confirmemos la política de una carrera, la tratamos como solo contacto: únicamente os ponemos en contacto, sin dinero a través de la plataforma.',
	'legal.terms.userHeading': 'Tus responsabilidades',
	'legal.terms.userBody':
		'Publica los dorsales con información veraz, respeta las normas de transferencia de cada carrera y nunca pidas más del valor nominal original del dorsal. No uses Bibseller para nada ilícito.',
	'legal.terms.conductHeading': 'Conducta prohibida e infractores reincidentes',
	'legal.terms.conductBody':
		'Podemos retirar anuncios o mensajes y restringir o cerrar cuentas que incumplan estos términos o las normas de una carrera. Las cuentas que infrinjan de forma reiterada se cerrarán.',
	'legal.terms.liabilityHeading': 'Responsabilidad',
	'legal.terms.liabilityBody':
		'Salvo cuando la plataforma procesa un pago de una carrera con reventa permitida, las transacciones son estrictamente entre usuarios y la plataforma no se responsabiliza de ellas. Nada de lo aquí dispuesto excluye la responsabilidad que no pueda excluirse según la ley aplicable.',
	'legal.terms.changesHeading': 'Cambios y contacto',
	'legal.terms.changesBody':
		'Podemos actualizar estos términos; los cambios importantes se publicarán en esta página. Las consultas, a través de nuestra página de contacto.',

	// Política de privacidad (M7.2, #54). BORRADOR pendiente de revisión legal (D10).
	'legal.privacy.title': 'Política de privacidad - Bibseller',
	'legal.privacy.metaDescription': 'Cómo trata Bibseller tus datos personales.',
	'legal.privacy.heading': 'Política de privacidad',
	'legal.privacy.intro':
		'Esta política explica qué datos personales tiene Bibseller, por qué, y los derechos que tienes sobre ellos.',
	'legal.privacy.dataHeading': 'Qué guardamos',
	'legal.privacy.dataBody':
		'Tu cuenta (correo, nombre visible, idioma preferido, país), los dorsales que publicas, los mensajes e imágenes que intercambias en el chat y las denuncias o bloqueos que realizas.',
	'legal.privacy.basisHeading': 'Por qué lo guardamos',
	'legal.privacy.basisBody':
		'Para gestionar el mercado que pediste usar: mostrar tus anuncios, entregar tus mensajes y mantener la plataforma segura. Nunca vendemos tus datos.',
	'legal.privacy.retentionHeading': 'Cuánto tiempo lo conservamos',
	'legal.privacy.retentionBody':
		'Los datos de la cuenta se conservan mientras tu cuenta esté activa. Los mensajes del chat se eliminan según un calendario de retención continuo. Algunos registros pueden conservarse más tiempo cuando la ley lo exige.',
	'legal.privacy.rightsHeading': 'Tus derechos',
	'legal.privacy.rightsBody':
		'Puedes acceder a tus datos personales, corregirlos, exportarlos o eliminarlos. La exportación y eliminación con autoservicio llegarán con nuestras herramientas completas de RGPD; mientras tanto, escríbenos y atenderemos tu solicitud.',
	'legal.privacy.cookiesHeading': 'Cookies',
	'legal.privacy.cookiesBody':
		'Usamos una única cookie estrictamente necesaria para mantener tu sesión iniciada. No usamos cookies de analítica ni de seguimiento, así que no hay banner de consentimiento.',
	'legal.privacy.contactHeading': 'Contacto',
	'legal.privacy.contactBody':
		'Para cualquier consulta o solicitud sobre privacidad, escríbenos a través de nuestra página de contacto.',

	// Contacto / aviso legal (M7.2, #54). Punto de contacto único a efectos de la DSA.
	'legal.contact.title': 'Contacto - Bibseller',
	'legal.contact.metaDescription':
		'Cómo contactar con Bibseller, incluido nuestro punto de contacto para las autoridades.',
	'legal.contact.heading': 'Contacto',
	'legal.contact.intro': 'Cómo ponerte en contacto con Bibseller.',
	'legal.contact.operatorHeading': 'Quién gestiona Bibseller',
	'legal.contact.operatorBody':
		'Bibseller se gestiona como un mercado sin ánimo de lucro. La entidad responsable, el domicilio y un correo de contacto se publicarán aquí antes de la beta pública.',
	'legal.contact.dsaHeading': 'Punto de contacto (Ley de Servicios Digitales)',
	'legal.contact.dsaBody':
		'También será nuestro punto de contacto único para usuarios y autoridades conforme a la Ley de Servicios Digitales de la UE. Puedes escribirnos en español o en inglés.',
	'legal.contact.reportHeading': 'Denunciar un anuncio',
	'legal.contact.reportBody':
		'Para denunciar un anuncio ilícito o que incumpla las normas, usa el botón Denunciar de ese anuncio: revisamos todas las denuncias.'
};

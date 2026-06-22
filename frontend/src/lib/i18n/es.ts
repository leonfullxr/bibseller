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
	'banner.suggestText': 'Esta página está disponible en español.',
	'banner.suggestAccept': 'Ver en español',
	'banner.suggestDismiss': 'Ahora no',
	'footer.tagline': 'Sin comisiones, en toda la UE. Sin ánimo de lucro por diseño.',
	'footer.github': 'GitHub',

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
	'error.backHome': 'Volver al inicio'
};

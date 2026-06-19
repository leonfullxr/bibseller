// The English message catalogue and the source of truth for the key set: every
// other locale is a Partial of these keys (es.ts), and `t()` falls back here for
// any key a translation has not filled yet. Flat dotted keys so `MessageKey`
// gives autocomplete and exhaustiveness for free - no nested-path type magic.
// `{name}`-style placeholders are filled by t(key, params). Grouped by area.

export const en = {
	// Language switcher + locale names.
	'lang.en': 'English',
	'lang.es': 'Español',
	'lang.switch': 'Language',

	// Top nav + global chrome.
	'nav.races': 'Races',
	'nav.sell': 'Sell',
	'nav.myListings': 'My listings',
	'nav.inbox': 'Inbox',
	'nav.login': 'Log in',
	'nav.register': 'Register',
	'nav.logout': 'Log out',
	'banner.verifyEmail': 'Verify your email to unlock selling and chat.',
	'banner.resend': 'Resend email',
	'footer.tagline': 'Zero commission, EU-wide. Non-profit by design.',
	'footer.github': 'GitHub',

	// Home.
	'home.title': 'Bibseller - race bibs find new runners',
	'home.metaDescription':
		"Non-profit, EU-wide marketplace connecting runners who can't start with runners who missed registration.",
	'home.heroTitle': 'Race bibs find',
	'home.heroTitleHighlight': 'new runners',
	'home.tagline':
		"Injured? Plans changed? Missed registration? A zero-commission marketplace that connects sellers and buyers of race bibs - always within each race's own rules.",
	'home.searchPlaceholder': 'Search a race or city…',
	'home.search': 'Search',
	'home.browseAll': 'or browse all races',
	'home.apiUnreachable': 'API unreachable - run',
	'home.upcoming': 'Upcoming races',
	'home.seeAll': 'See all',
	'home.modePlatformSaleName': 'Platform sale',
	'home.modePlatformSaleDesc':
		'The race allows resale: list, chat, and pay securely through the platform.',
	'home.modeOfficialName': 'Official process',
	'home.modeOfficialDesc':
		'The race runs its own name change: we connect you and link the official procedure.',
	'home.modeConnectName': 'Connect only',
	'home.modeConnectDesc':
		'Restricted or unverified races: we provide the chat, the rest stays between you two.',
	'home.underConstruction': 'Under construction - follow the',
	'home.roadmap': 'roadmap',

	// Policy words (moved from $lib/policy.ts - the facts stay there, words live here).
	'policy.label.platform_sale': 'Resale allowed',
	'policy.label.official_only': 'Official transfer',
	'policy.label.connect_only': 'Chat only',
	'policy.label.unknown': 'Policy unverified',
	'policy.disclaimer.platform_sale.title': 'This race allows bib resale.',
	'policy.disclaimer.platform_sale.body':
		'Agree with the seller in chat, then pay securely through the platform - funds are held until the transfer is confirmed. Zero commission.',
	'policy.disclaimer.official_only.title': 'This race runs its own official name-change process.',
	'policy.disclaimer.official_only.body':
		'Find each other and agree on the details here - the transfer itself (and any official fee) goes through the race organizer. The platform never handles money for this race.',
	'policy.disclaimer.connect_only.title': 'This race restricts bib transfers.',
	'policy.disclaimer.connect_only.body':
		"The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything.",
	'policy.disclaimer.unknown.title':
		'Transfer policy not verified yet - treat this race as chat-only.',
	'policy.disclaimer.unknown.body':
		"The platform only connects you: it handles no money here and takes no responsibility for any arrangement between you and the other party. The race's own rules apply - check them before agreeing to anything.",
	'policy.officialLink': 'Official transfer process',

	// Sport filter values (the option value stays the enum; only the label is translated).
	'sport.running': 'Running',
	'sport.trail': 'Trail',
	'sport.triathlon': 'Triathlon',
	'sport.cycling': 'Cycling',
	'sport.obstacle': 'Obstacle',
	'sport.other': 'Other',

	// Races browse + filters.
	'races.title': 'Browse races - Bibseller',
	'races.metaDescription': 'Find race bibs for sale across EU running events.',
	'races.heading': 'Browse races',
	'races.filter.search': 'Search',
	'races.filter.searchPlaceholder': 'Race or city…',
	'races.filter.country': 'Country',
	'races.filter.sport': 'Sport',
	'races.filter.policy': 'Transfer policy',
	'races.filter.all': 'All',
	'races.filter.submit': 'Filter',
	'races.empty': 'No races match those filters.',
	'races.clearFilters': 'Clear filters',
	'races.nextPage': 'Next page ->',

	// Race + listing cards. CLDR plural forms via Intl.PluralRules (no ICU lib, D14);
	// en/es use one|other - add .few/.many keys if a locale needs them.
	'raceCard.bibs.one': '{n} bib listed',
	'raceCard.bibs.other': '{n} bibs listed',
	'listingCard.priceOnRequest': 'Price on request',
	'listingCard.belowFace': 'below face value',
	'listingCard.listedBy': 'Listed by {name}',

	// Race detail.
	'raceDetail.title': '{name} - bibs for sale - Bibseller',
	'raceDetail.metaDescription': 'Bibs for {name} ({date}, {city}).',
	'raceDetail.back': 'Back to all races',
	'raceDetail.website': 'Race website',
	'raceDetail.bibsForSale.one': '{n} bib for sale',
	'raceDetail.bibsForSale.other': '{n} bibs for sale',
	'raceDetail.sellCta': 'Sell your bib',
	'raceDetail.empty': 'No bibs listed for this race yet.',
	'raceDetail.emptyHint': 'Selling yours? Listing opens soon.',

	// Listing detail.
	'listingDetail.title': 'Bib for {name} - Bibseller',
	'listingDetail.back': 'Back to {name}',
	'listingDetail.heading': 'Bib for {name}',
	'listingDetail.unavailable': 'This listing is no longer available ({status}).',
	'listingDetail.listedByOn': 'Listed by {name} on {date}',
	'listingDetail.contact': 'Contact the seller',
	'listingDetail.toMessageSeller': 'to message the seller.',
	'listingDetail.verifyToMessage': 'Verify your email to message the seller.',
	'listingDetail.accountSettings': 'Account settings',
	'listingDetail.ownPre': 'This is your listing - manage it from',
	'listingDetail.yourListings': 'your listings',
	'listingDetail.messageAria': 'Message to the seller',
	'listingDetail.messagePlaceholder': 'Hi - is this bib still available?',
	'listingDetail.ackText':
		"I understand the platform handles no money and takes no responsibility for this transfer - the race's own rules apply.",
	'listingDetail.send': 'Send message',

	// Report a listing / content.
	'report.summary': 'Report this listing',
	'report.reason.forbidden_transfer': 'Forbidden transfer',
	'report.reason.scam': 'Scam',
	'report.reason.offensive': 'Offensive',
	'report.reason.other': 'Other',
	'report.reasonAria': 'Reason for report',
	'report.detailsAria': 'Report details (optional)',
	'report.detailsPlaceholder': 'Details (optional)',
	'report.success': 'Thanks - this listing has been reported.',
	'report.failed': 'Could not file the report. Try again.',
	'report.networkError': 'Network error - try again.',
	'report.submitting': 'Reporting...',
	'report.submit': 'Submit report',

	// Listing CTA (the buy path is an honest disabled stub until M6).
	'listingCta.buy': 'Buy securely - coming soon',
	'listingCta.buyTitle': 'Secure checkout arrives with payments (M6)',

	// Inbox (thread list).
	'inbox.title': 'Inbox - Bibseller',
	'inbox.heading': 'Inbox',
	'inbox.emptyPre': 'No conversations yet. Browse',
	'inbox.emptyRacesLink': 'races',
	'inbox.emptyPost': 'and contact a seller to start one.',
	'role.seller': 'seller',
	'role.buyer': 'buyer',

	// Chat thread.
	'chat.title': 'Chat with {name} - Bibseller',
	'chat.back': 'Back to inbox',
	'chat.about': 'about',
	'chat.block': 'Block',
	'chat.unblock': 'Unblock',
	'chat.sharedImage': 'Shared image',
	'chat.reportMsg': 'report',
	'chat.messageAria': 'Your message',
	'chat.messagePlaceholder': 'Write a message, or attach an image...',
	'chat.attachAria': 'Attach an image (JPEG or PNG)',
	'chat.send': 'Send',
	'chat.sending': 'Sending...',
	'chat.imageTooLarge': 'That image is too large (5 MB max).',
	'chat.tooFast': 'You are sending messages too fast - wait a moment.',
	'chat.sendFailed': 'Could not send your message. Try again.',
	'chat.networkError': 'Network error - check your connection.',
	'chat.blockConfirm': 'Block {name}? Neither of you will be able to message the other.',
	'chat.blocked': 'User blocked.',
	'chat.blockFailed': 'Could not block the user.',
	'chat.unblocked': 'User unblocked.',
	'chat.unblockFailed': 'Could not unblock the user.',
	'chat.networkRetry': 'Network error - try again.',
	'chat.reportConfirm': 'Report this message to the moderators?',
	'chat.messageReported': 'Message reported.',
	'chat.messageReportFailed': 'Could not report the message.',

	// Sell - race search.
	'sell.title': 'Sell a bib - Bibseller',
	'sell.heading': 'Sell a bib',
	'sell.lede': 'Find your race, then list your bib. You set the price (capped at face value).',
	'sell.verifyNotice':
		'Verify your email to publish a listing - you can still find your race below.',
	'sell.searchAria': 'Search races by name or city',
	'sell.emptyPre': 'No upcoming races match. Try another search, or',
	'sell.browseAllLink': 'browse all races',
	'sell.sellHere': 'Sell here',

	// Sell - listing form.
	'sellForm.title': 'List your bib for {name} - Bibseller',
	'sellForm.back': 'Back to race search',
	'sellForm.heading': 'List your bib',
	'sellForm.verifyNotice': 'Verify your email to publish a listing.',
	'sellForm.publish': 'Publish listing',

	// Shared listing form fields (create + edit).
	'listingFields.price': 'Asking price (EUR)',
	'listingFields.pricePlaceholder': 'e.g. 45',
	'listingFields.original': 'Original price / face value (EUR)',
	'listingFields.optional': 'optional',
	'listingFields.hint': 'Enter the face value and your asking price is capped at it - no scalping.',
	'listingFields.description': 'Description',
	'listingFields.descriptionPlaceholder': 'optional - size, pickup details, etc.',

	// My listings.
	'myListings.title': 'My listings - Bibseller',
	'myListings.heading': 'My listings',
	'myListings.emptyPre': 'You have no listings yet.',
	'myListings.listABib': 'List a bib',
	'myListings.edit': 'Edit',
	'myListings.cancel': 'Cancel',

	// Edit listing.
	'editListing.title': 'Edit listing - Bibseller',
	'editListing.back': 'Back to my listings',
	'editListing.heading': 'Edit listing',
	'editListing.save': 'Save changes',

	// Auth - shared field labels.
	'auth.email': 'Email',
	'auth.password': 'Password',

	// Log in.
	'login.title': 'Log in - Bibseller',
	'login.forgot': 'Forgot your password?',
	'login.newHere': 'New here?',
	'login.createAccount': 'Create an account',

	// Register.
	'register.title': 'Create account - Bibseller',
	'register.heading': 'Create account',
	'register.displayName': 'Display name',
	'register.haveAccount': 'Already have an account?',

	// Forgot password (request a reset link).
	'forgot.title': 'Reset password - Bibseller',
	'forgot.heading': 'Reset your password',
	'forgot.lede': "Enter your email and we'll send you a reset link.",
	'forgot.sent':
		"If an account exists for that address, we've sent a link to reset your password. Check your inbox.",
	'forgot.submit': 'Send reset link',
	'forgot.backToLogin': 'Back to log in',

	// Reset password (set a new one).
	'reset.title': 'Set a new password - Bibseller',
	'reset.heading': 'Set a new password',
	'reset.done':
		"Your password has been updated. You've been signed out everywhere - sign in with your new password.",
	'reset.missingToken': 'This reset link is missing its token. Request a new one.',
	'reset.requestLink': 'Request a reset link',
	'reset.newPassword': 'New password',
	'reset.confirmPassword': 'Confirm password',
	'reset.submit': 'Update password',

	// Verify email (landing page).
	'verify.title': 'Verify email - Bibseller',
	'verify.okHeading': 'Email verified',
	'verify.okBody': "Your email address is confirmed - you're all set.",
	'verify.continue': 'Continue',
	'verify.invalidHeading': 'Link invalid or expired',
	'verify.invalidBody':
		'This verification link is no longer valid. Sign in and request a fresh one.',
	'verify.signIn': 'Sign in',
	'verify.missingHeading': 'Nothing to verify',
	'verify.missingBody': 'Open the verification link from your email to confirm your address.',
	'verify.home': 'Home',
	'verify.errorHeading': 'Something went wrong',
	'verify.errorBody': "We couldn't verify your email right now. Please try again in a moment.",

	// Settings.
	'settings.title': 'Settings - Bibseller',
	'settings.heading': 'Settings',
	'settings.profile': 'Profile',
	'settings.country': 'Country',
	'settings.countryNotSet': 'Not set',
	'settings.profileUpdated': 'Profile updated.',
	'settings.save': 'Save',
	'settings.password': 'Password',
	'settings.currentPassword': 'Current password',
	'settings.confirmNewPassword': 'Confirm new password',
	'settings.passwordChanged': 'Password changed. Other devices have been signed out.',
	'settings.changePassword': 'Change password',
	'settings.sessions': 'Sessions',
	'settings.sessionsNote': 'Sign out of Bibseller on every device, including this one.',
	'settings.logoutAll': 'Log out all devices',
	'settings.deleteAccount': 'Delete account',
	'settings.deleteNote':
		'Permanently delete your account and its data. Available once full GDPR tooling ships (M7).',
	'settings.deleteTitle': 'Account deletion arrives with trust and safety (M7)',
	'settings.deleteSoon': 'Delete account - coming soon',

	// Error page.
	'error.notFound': "That page doesn't exist.",
	'error.generic': 'Something went wrong.',
	'error.backHome': 'Back home'
} as const;

export type MessageKey = keyof typeof en;

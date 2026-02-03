<?php

declare(strict_types=1);

return [
    'secret' => env('JWT_SECRET', 'change-me'),
    'issuer' => env('JWT_ISSUER', 'task-system'),
    'ttl' => (int) env('JWT_TTL', 7200),
];

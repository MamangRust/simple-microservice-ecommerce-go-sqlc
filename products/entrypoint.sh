#!/bin/sh
echo "📦 Running migrations..."
if ./migrate up; then
    echo "✅ Migration completed successfully"
else
    EXIT_CODE=$?
    echo "❌ Migration failed with exit code $EXIT_CODE"
    exit $EXIT_CODE
fi

echo ""
echo "🚀 Starting product service..."
echo "========================================="
exec ./product
// Custom validation rules for form inputs
import { defineRule, configure } from 'vee-validate'
import { required, email, min, max, numeric } from '@vee-validate/rules'

// Configure VeeValidate with custom error messages
export function setupValidation() {
  // Define built-in rules
  defineRule('required', required)
  defineRule('email', email)
  defineRule('min', min)
  defineRule('max', max)
  defineRule('numeric', numeric)
  
  // Custom rule for minimum age validation
  defineRule('minAge', (value, [limit]) => {
    if (!value || value < limit) {
      return `Age must be at least ${limit} years`
    }
    return true
  })
  
  // Custom rule for unique email (async validation)
  defineRule('uniqueEmail', async (value, _, ctx) => {
    if (!value) return true
    
    try {
      // Check if email exists in database
      const response = await api.checkEmailExists(value)
      
      // Skip validation if editing same user
      if (ctx.form.userId && response.userId === ctx.form.userId) {
        return true
      }
      
      if (response.exists) {
        return 'Email already exists'
      }
      
      return true
    } catch {
      // Allow validation to pass on API error
      return true
    }
  })
  
  // Custom rule for password strength
  defineRule('strongPassword', (value) => {
    if (!value) return true
    
    const hasUpperCase = /[A-Z]/.test(value)
    const hasLowerCase = /[a-z]/.test(value)
    const hasNumbers = /\d/.test(value)
    const hasNonalphas = /\W/.test(value)
    
    if (!hasUpperCase || !hasLowerCase || !hasNumbers || !hasNonalphas) {
      return 'Password must contain uppercase, lowercase, number and special character'
    }
    
    return true
  })
  
  // Configure global validation messages
  configure({
    generateMessage: (context) => {
      const messages = {
        required: `${context.field} is required`,
        email: `${context.field} must be a valid email`,
        min: `${context.field} must be at least ${context.rule.params[0]} characters`,
        max: `${context.field} must not exceed ${context.rule.params[0]} characters`,
        numeric: `${context.field} must be a number`,
        minAge: `User must be at least ${context.rule.params[0]} years old`,
        uniqueEmail: 'This email is already registered',
        strongPassword: 'Please choose a stronger password'
      }
      
      return messages[context.rule.name] || `${context.field} is invalid`
    },
    
    validateOnBlur: true,
    validateOnChange: false,
    validateOnInput: false,
    validateOnModelUpdate: true
  })
}

// Validation schemas for different forms
export const validationSchemas = {
  // User creation validation schema
  createUser: {
    name: {
      required: true,
      min: 2,
      max: 100
    },
    email: {
      required: true,
      email: true,
      uniqueEmail: true
    },
    age: {
      required: true,
      numeric: true,
      minAge: 18
    },
    password: {
      required: true,
      min: 6,
      strongPassword: false // Optional strong password
    },
    role: {
      required: true
    }
  },
  
  // User update validation schema
  updateUser: {
    name: {
      required: true,
      min: 2,
      max: 100
    },
    email: {
      required: true,
      email: true
    },
    age: {
      required: true,
      numeric: true,
      minAge: 18
    },
    role: {
      required: true
    }
  },
  
  // Login validation schema
  login: {
    email: {
      required: true,
      email: true
    },
    password: {
      required: true,
      min: 6
    }
  }
}

// Helper function to validate single field
export async function validateField(value, rules) {
  for (const rule of Object.keys(rules)) {
    const ruleValue = rules[rule]
    const validationResult = await defineRule[rule](value, ruleValue)
    
    if (validationResult !== true) {
      return validationResult
    }
  }
  
  return true
}

// Export validation utilities
export default {
  setupValidation,
  validationSchemas,
  validateField
}
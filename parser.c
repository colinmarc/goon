#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "parser.h"

gn_parser_context_t *gn_context;

static void gn_init_context(gn_parser_context_t *context) {
    context->position = 0;
    context->next_symbol = 0;
    context->buffer_len = 0;
    context->buffer_offset = 0;
    context->error = 0;
}

gn_parser_context_t *gn_global_context() {
    if (gn_context == NULL) {
        gn_context = malloc(sizeof(*gn_context));
        gn_init_context(gn_context);
    }

    return gn_context;
}

static gn_ast_node_t *gn_create_node(gn_ast_node_type_t node_type, int value) {
    gn_ast_node_t *node = malloc(sizeof(*node));

    node->node_type = node_type;
    node->value = value;
    node->num_children = 0;

    return node;
}

static gn_ast_node_t *gn_top(gn_parser_context_t *context) {
    return context->stack[context->position - 1];
}

static void gn_push(gn_parser_context_t *context, gn_ast_node_t *node)	{
    //printf("pushing %d/%d/%d\n", node->node_type, node->value, node->num_children);
    context->stack[context->position] = node;
    context->position++;
}

static gn_ast_node_t *gn_pop(gn_parser_context_t *context) {
    if (context->position <= 0) return NULL;

    gn_ast_node_t *node = gn_top(context);
    context->position--;
    return node;
}

static void gn_reduce(gn_parser_context_t *context,
                      gn_ast_node_type_t node_type, int num_children) {
    //printf("reducing! %d, %d\n", node_type, num_children);
    gn_ast_node_t **children = malloc(num_children * sizeof(gn_ast_node_t *));

    int i;
    for(i = num_children - 1; i >= 0; i--) {
        gn_ast_node_t *node = gn_pop(context);
        children[i] = node;
    }

    gn_ast_node_t *node = gn_create_node(node_type, -1);
    node->children = children;
    node->num_children = num_children;

    gn_push(context, node);
}

static void gn_number_node(gn_parser_context_t *context, char *text) {
    int value = atoi(text);
    gn_push(context, gn_create_node(GN_AST_NUMBER, value));
}

static void gn_symbol_node(gn_parser_context_t *context, char *text) {
    int i = context->next_symbol++;
    printf("saving symbol %d: %s\n", i, text);
    context->symbols[i] = text;
    printf("symbols[0]: %s\n", context->symbols[0]);

    gn_push(context, gn_create_node(GN_AST_SYMBOL, i));
}

static int gn_read_input(char *buffer, int *result, int max_size) {
    gn_parser_context_t *context = gn_global_context();

    if(context->buffer_len == 0) {
        *result = 0;
        return 0;
    }

    int num_bytes = context->buffer_len;
    if (num_bytes > max_size) num_bytes = max_size;

    int i;
    for (i = 0; i < num_bytes; i++) {
        buffer[i] = context->input_buffer[context->buffer_offset+i];
    }

    *result = num_bytes;
    context->buffer_offset += num_bytes;
    context->buffer_len -= num_bytes;

    return 0;
}

static void gn_feed_input(gn_parser_context_t *context,
                          char *new_input, int new_input_len) {
    context->input_buffer = new_input;
    context->buffer_len = new_input_len;
    context->buffer_offset = 0;
}

#undef YY_INPUT
#define YY_INPUT(buf, result, max_size) gn_read_input(buf, &result, max_size);

#include "goon.peg.c"

// public

gn_ast_node_t *gn_parse(gn_parser_context_t *context, char *buffer, int len) {
    gn_init_context(context);

    gn_feed_input(context, buffer, len);
    yyparse();

    if (context->error) {
        return NULL;
    }

    gn_ast_node_t *root = gn_pop(context);
    return root;
}

gn_ast_node_t *gn_child_at(gn_ast_node_t *parent, int idx) {
    return parent->children[idx];
}

char *gn_get_symbol(gn_parser_context_t *context, gn_ast_node_t *node) {
    //printf("getting symbol %d: %s\n", node->value, context->symbols[node->value]);
    printf("symbols[0]: %s\n", context->symbols[0]);
    return context->symbols[node->value];
}

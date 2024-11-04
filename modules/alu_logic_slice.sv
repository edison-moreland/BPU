module alu_logic_slice (
    input  logic a, b,
    input  logic prev_eq, prev_greater,
    output logic result_xor,
    output logic result_or,
    output logic result_and,
    output logic result_not,
    output logic is_eq, is_greater
);

    always @* begin
        is_eq <= (a == b) & prev_eq;
        is_greater <= ((a > b) & prev_eq) | prev_greater;

        result_xor <= a ^ b;
        result_or <= a | b;
        result_and <= a & b;
        result_not <= !a;
    end

endmodule